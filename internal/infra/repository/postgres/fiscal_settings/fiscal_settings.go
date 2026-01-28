package fiscalsettingsrepositorybun

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	fiscalsettingsentity "github.com/willjrcom/sales-backend-go/internal/domain/fiscal_settings"
)

type FiscalSettings struct {
	bun.BaseModel `bun:"table:fiscal_settings,alias:fs"`

	ID                 uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	CompanyID          uuid.UUID `bun:"company_id,type:uuid,notnull"`
	IsActive           bool      `bun:"is_active"`
	InscricaoEstadual  string    `bun:"inscricao_estadual"`
	RegimeTributario   int       `bun:"regime_tributario"`
	CNAE               string    `bun:"cnae"`
	CRT                int       `bun:"crt"`
	SimplesNacional    bool      `bun:"simples_nacional"`
	InscricaoMunicipal string    `bun:"inscricao_municipal"`

	// Preferences
	DiscriminaImpostos      bool `bun:"discrimina_impostos"`
	EnviarEmailDestinatario bool `bun:"enviar_email_destinatario"`

	// Company Identity
	BusinessName string `bun:"business_name"`
	TradeName    string `bun:"trade_name"`
	Cnpj         string `bun:"cnpj"`
	Email        string `bun:"email"`
	Phone        string `bun:"phone"`

	// Address Fields (Flattened for now, or JSONB? Flattened usually better for simple queries, but Address struct in domain suggests decoupling. Let's start flattened as per company address pattern in some projects, or JSONB if flexible. The DTO had flattened fields. Let's flatten here to match.)
	Street       string `bun:"street"`
	Number       string `bun:"number"`
	Complement   string `bun:"complement"`
	Neighborhood string `bun:"neighborhood"`
	City         string `bun:"city"`
	UF           string `bun:"uf"`
	Cep          string `bun:"cep"`

	CreatedAt time.Time  `bun:"created_at,default:current_timestamp"`
	UpdatedAt *time.Time `bun:"updated_at"`
}

func (m *FiscalSettings) ToDomain() *fiscalsettingsentity.FiscalSettings {
	return &fiscalsettingsentity.FiscalSettings{
		Entity: entity.Entity{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: *m.UpdatedAt,
		},
		CompanyID:               m.CompanyID,
		IsActive:                m.IsActive,
		InscricaoEstadual:       m.InscricaoEstadual,
		RegimeTributario:        m.RegimeTributario,
		CNAE:                    m.CNAE,
		CRT:                     m.CRT,
		SimplesNacional:         m.SimplesNacional,
		InscricaoMunicipal:      m.InscricaoMunicipal,
		DiscriminaImpostos:      m.DiscriminaImpostos,
		EnviarEmailDestinatario: m.EnviarEmailDestinatario,
		BusinessName:            m.BusinessName,
		TradeName:               m.TradeName,
		Cnpj:                    m.Cnpj,
		Email:                   m.Email,
		Phone:                   m.Phone,
		Address: fiscalsettingsentity.FiscalAddress{
			Street:       m.Street,
			Number:       m.Number,
			Complement:   m.Complement,
			Neighborhood: m.Neighborhood,
			City:         m.City,
			UF:           m.UF,
			Cep:          m.Cep,
		},
	}
}

func (m *FiscalSettings) FromDomain(d *fiscalsettingsentity.FiscalSettings) {
	m.ID = d.ID
	m.CompanyID = d.CompanyID
	m.IsActive = d.IsActive
	m.InscricaoEstadual = d.InscricaoEstadual
	m.RegimeTributario = d.RegimeTributario
	m.CNAE = d.CNAE
	m.CRT = d.CRT
	m.SimplesNacional = d.SimplesNacional
	m.InscricaoMunicipal = d.InscricaoMunicipal
	m.DiscriminaImpostos = d.DiscriminaImpostos
	m.EnviarEmailDestinatario = d.EnviarEmailDestinatario
	m.BusinessName = d.BusinessName
	m.TradeName = d.TradeName
	m.Cnpj = d.Cnpj
	m.Email = d.Email
	m.Phone = d.Phone
	m.Street = d.Address.Street
	m.Number = d.Address.Number
	m.Complement = d.Address.Complement
	m.Neighborhood = d.Address.Neighborhood
	m.City = d.Address.City
	m.UF = d.Address.UF
	m.Cep = d.Address.Cep
	m.CreatedAt = d.CreatedAt
	m.UpdatedAt = &d.UpdatedAt
}

type FiscalSettingsRepositoryBun struct {
	db *bun.DB
}

func NewFiscalSettingsRepositoryBun(db *bun.DB) *FiscalSettingsRepositoryBun {
	return &FiscalSettingsRepositoryBun{db: db}
}

func (r *FiscalSettingsRepositoryBun) Create(ctx context.Context, settings *fiscalsettingsentity.FiscalSettings) error {
	model := &FiscalSettings{}
	model.FromDomain(settings)
	_, err := r.db.NewInsert().Model(model).Exec(ctx)
	return err
}

func (r *FiscalSettingsRepositoryBun) Update(ctx context.Context, settings *fiscalsettingsentity.FiscalSettings) error {
	model := &FiscalSettings{}
	model.FromDomain(settings)
	_, err := r.db.NewUpdate().Model(model).WherePK().Exec(ctx)
	return err
}

func (r *FiscalSettingsRepositoryBun) GetByCompanyID(ctx context.Context, companyID uuid.UUID) (*fiscalsettingsentity.FiscalSettings, error) {
	model := &FiscalSettings{}
	err := r.db.NewSelect().Model(model).Where("company_id = ?", companyID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}
