package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"nestars-bot/internal/domain"
)

const (
	usersColl    = "users"
	profilesColl = "profiles"
	ordersColl   = "orders"
	countersColl = "counters"
)

// mongo-структуры хранения.

type mongoUser struct {
	ID           string    `bson:"_id"`
	TelegramID   int64     `bson:"telegram_id"`
	Username     string    `bson:"username"`
	FirstName    string    `bson:"first_name"`
	LastName     string    `bson:"last_name"`
	LanguageCode string    `bson:"language_code"`
	ReferralCode string    `bson:"referral_code"`
	ReferredBy   *string   `bson:"referred_by,omitempty"`
	CreatedAt    time.Time `bson:"created_at"`
}

type mongoProfile struct {
	ID                 string  `bson:"_id"`
	UserID             string  `bson:"user_id"`
	TelegramID         int64   `bson:"telegram_id"`
	ReferralCode       string  `bson:"referral_code"`
	BalanceUSD         float64 `bson:"balance_usd"`
	BalanceTON         float64 `bson:"balance_ton"`
	BoughtStars        int     `bson:"bought_stars"`
	BonusStars         int     `bson:"bonus_stars"`
	ReferralBonusStars int     `bson:"referral_bonus_stars"`
	ReferralCount      int     `bson:"referral_count"`
}

type mongoOrder struct {
	ID                string    `bson:"_id"`
	OrderNumber       int       `bson:"order_number"`
	UserID            string    `bson:"user_id"`
	RecipientUsername string    `bson:"recipient_username"`
	ProductType       string    `bson:"product_type"`
	ProductID         string    `bson:"product_id"`
	StarsAmount       *int      `bson:"stars_amount,omitempty"`
	Amount            float64   `bson:"amount"`
	Currency          string    `bson:"currency"`
	Status            string    `bson:"status"`
	PaymentMethod     string    `bson:"payment_method,omitempty"`
	PaymentURL        string    `bson:"payment_url,omitempty"`
	CreatedAt         time.Time `bson:"created_at"`
}

// userRepo — реализация domain.UserRepository на MongoDB.
type userRepo struct {
	db *mongo.Database
}

func NewUserRepo(db *mongo.Database) domain.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) CreateOrGet(ctx context.Context, tg *domain.TelegramUser, refCode string) (*domain.User, error) {
	coll := r.db.Collection(usersColl)

	var u mongoUser
	err := coll.FindOne(ctx, bson.M{"telegram_id": tg.ID}).Decode(&u)
	if err == nil {
		return toDomainUser(u), nil
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	id := newID()
	u = mongoUser{
		ID:           id,
		TelegramID:   tg.ID,
		Username:     tg.Username,
		FirstName:    tg.FirstName,
		LastName:     tg.LastName,
		LanguageCode: tg.LanguageCode,
		ReferralCode: generateRefCode(id),
		CreatedAt:    time.Now().UTC(),
	}
	if refCode != "" {
		u.ReferredBy = &refCode
	}

	if _, err := coll.InsertOne(ctx, u); err != nil {
		return nil, err
	}

	p := mongoProfile{
		ID:           newID(),
		UserID:       u.ID,
		TelegramID:   u.TelegramID,
		ReferralCode: u.ReferralCode,
	}
	if _, err := r.db.Collection(profilesColl).InsertOne(ctx, p); err != nil {
		return nil, err
	}

	return toDomainUser(u), nil
}

func (r *userRepo) GetUserWithProfile(ctx context.Context, telegramID int64) (*domain.User, *domain.Profile, error) {
	var u mongoUser
	err := r.db.Collection(usersColl).FindOne(ctx, bson.M{"telegram_id": telegramID}).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil, fmt.Errorf("user not found")
		}
		return nil, nil, err
	}

	var p mongoProfile
	err = r.db.Collection(profilesColl).FindOne(ctx, bson.M{"telegram_id": telegramID}).Decode(&p)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// Если профиль потерялся — создаём пустой.
			p = mongoProfile{
				ID:         newID(),
				UserID:     u.ID,
				TelegramID: u.TelegramID,
			}
			if _, e := r.db.Collection(profilesColl).InsertOne(ctx, p); e != nil {
				return nil, nil, e
			}
		} else {
			return nil, nil, err
		}
	}

	return toDomainUser(u), toDomainProfile(p), nil
}

func (r *userRepo) GetProfile(ctx context.Context, userID string) (*domain.Profile, error) {
	var p mongoProfile
	err := r.db.Collection(profilesColl).FindOne(ctx, bson.M{"user_id": userID}).Decode(&p)
	if err != nil {
		return nil, err
	}
	return toDomainProfile(p), nil
}

// orderRepo — реализация domain.OrderRepository на MongoDB.
type orderRepo struct {
	db *mongo.Database
}

func NewOrderRepo(db *mongo.Database) domain.OrderRepository {
	return &orderRepo{db: db}
}

func (r *orderRepo) Create(ctx context.Context, order *domain.Order) error {
	mo := mongoOrder{
		ID:                order.ID,
		UserID:            order.UserID,
		RecipientUsername: order.RecipientUsername,
		ProductType:       order.ProductType,
		ProductID:         order.ProductID,
		StarsAmount:       order.StarsAmount,
		Amount:            order.Amount,
		Currency:          order.Currency,
		Status:            order.Status,
		PaymentMethod:     order.PaymentMethod,
		PaymentURL:        order.PaymentURL,
		CreatedAt:         time.Now().UTC(),
	}
	if mo.ID == "" {
		mo.ID = newID()
	}
	if mo.OrderNumber == 0 {
		n, err := nextOrderNumber(ctx, r.db)
		if err != nil {
			return err
		}
		mo.OrderNumber = n
	}
	if mo.Status == "" {
		mo.Status = "pending"
	}

	_, err := r.db.Collection(ordersColl).InsertOne(ctx, mo)
	if err != nil {
		return err
	}

	order.ID = mo.ID
	order.OrderNumber = mo.OrderNumber
	order.Status = mo.Status
	return nil
}

func (r *orderRepo) UpdatePayment(ctx context.Context, orderID, method string) error {
	_, err := r.db.Collection(ordersColl).UpdateOne(
		ctx,
		bson.M{"_id": orderID},
		bson.M{"$set": bson.M{"payment_method": method}},
	)
	return err
}

func (r *orderRepo) UpdatePaymentURL(ctx context.Context, orderID, url string) error {
	_, err := r.db.Collection(ordersColl).UpdateOne(
		ctx,
		bson.M{"_id": orderID},
		bson.M{"$set": bson.M{"payment_url": url}},
	)
	return err
}

func (r *orderRepo) FindByUserID(ctx context.Context, userID string) ([]domain.Order, error) {
	cursor, err := r.db.Collection(ordersColl).Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	out := []domain.Order{}
	for cursor.Next(ctx) {
		var mo mongoOrder
		if err := cursor.Decode(&mo); err != nil {
			continue
		}
		out = append(out, toDomainOrder(mo))
	}
	return out, cursor.Err()
}

func toDomainOrder(mo mongoOrder) domain.Order {
	o := domain.Order{
		ID:                mo.ID,
		OrderNumber:       mo.OrderNumber,
		UserID:            mo.UserID,
		RecipientUsername: mo.RecipientUsername,
		ProductType:       mo.ProductType,
		ProductID:         mo.ProductID,
		Amount:            mo.Amount,
		Currency:          mo.Currency,
		Status:            mo.Status,
		PaymentMethod:     mo.PaymentMethod,
		PaymentURL:        mo.PaymentURL,
	}
	if mo.StarsAmount != nil {
		o.StarsAmount = mo.StarsAmount
	}
	return o
}

func (r *orderRepo) ProcessBalancePayment(ctx context.Context, orderID string) error {
	orders := r.db.Collection(ordersColl)
	profiles := r.db.Collection(profilesColl)

	var o mongoOrder
	if err := orders.FindOne(ctx, bson.M{"_id": orderID}).Decode(&o); err != nil {
		return err
	}

	var p mongoProfile
	if err := profiles.FindOne(ctx, bson.M{"user_id": o.UserID}).Decode(&p); err != nil {
		return err
	}

	if p.BalanceUSD < o.Amount {
		return fmt.Errorf("insufficient balance")
	}

	inc := bson.M{"balance_usd": -o.Amount}
	if o.StarsAmount != nil {
		inc["bought_stars"] = *o.StarsAmount
	}

	if _, err := profiles.UpdateOne(ctx, bson.M{"_id": p.ID}, bson.M{"$inc": inc}); err != nil {
		return err
	}

	_, err := orders.UpdateOne(
		ctx,
		bson.M{"_id": orderID},
		bson.M{"$set": bson.M{"status": "paid", "payment_method": "Balance"}},
	)
	return err
}

// nextOrderNumber возвращает следующий порядковый номер заказа.
func nextOrderNumber(ctx context.Context, db *mongo.Database) (int, error) {
	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var doc struct {
		Seq int `bson:"seq"`
	}
	err := db.Collection(countersColl).FindOneAndUpdate(
		ctx,
		bson.M{"_id": "orders"},
		bson.M{"$inc": bson.M{"seq": 1}},
		opts,
	).Decode(&doc)
	if err != nil {
		return 0, err
	}
	return doc.Seq, nil
}

// EnsureIndexes создаёт необходимые индексы при старте приложения.
func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	indexes := []struct {
		coll string
		key  string
		uniq bool
	}{
		{usersColl, "telegram_id", true},
		{profilesColl, "user_id", true},
		{profilesColl, "telegram_id", false},
		{ordersColl, "order_number", true},
	}

	for _, idx := range indexes {
		opts := options.Index()
		if idx.uniq {
			opts = opts.SetUnique(true)
		}
		// Явное имя исключает конфликты автогенерации имени.
		opts = opts.SetName(fmt.Sprintf("%s_%s_idx", idx.coll, idx.key))

		_, err := db.Collection(idx.coll).Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys:    bson.D{{Key: idx.key, Value: 1}},
			Options: opts,
		})
		if err != nil {
			// Игнорируем ошибку, если индекс уже существует.
			var cmdErr mongo.CommandError
			if errors.As(err, &cmdErr) && (cmdErr.Code == 86 || cmdErr.Code == 68) {
				continue
			}
			return err
		}
	}

	return nil
}

func newID() string {
	return bson.NewObjectID().Hex()
}

func generateRefCode(id string) string {
	if len(id) >= 8 {
		return id[:8]
	}
	return id
}

func toDomainUser(u mongoUser) *domain.User {
	return &domain.User{
		ID:           u.ID,
		TelegramID:   u.TelegramID,
		Username:     u.Username,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		LanguageCode: u.LanguageCode,
		ReferralCode: u.ReferralCode,
		ReferredBy:   u.ReferredBy,
		CreatedAt:    u.CreatedAt,
	}
}

func toDomainProfile(p mongoProfile) *domain.Profile {
	return &domain.Profile{
		ID:                 p.ID,
		UserID:             p.UserID,
		TelegramID:         p.TelegramID,
		ReferralCode:       p.ReferralCode,
		BalanceUSD:         p.BalanceUSD,
		BalanceTON:         p.BalanceTON,
		BoughtStars:        p.BoughtStars,
		BonusStars:         p.BonusStars,
		ReferralBonusStars: p.ReferralBonusStars,
		ReferralCount:      p.ReferralCount,
	}
}
