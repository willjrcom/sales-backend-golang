package fiscalinvoice

// Fiscal constants for food/restaurant items
// These are default values for alimentação preparada (prepared food)
const (
	// NCM - Nomenclatura Comum do Mercosul
	// 21069090 - Preparações alimentícias não especificadas
	DefaultFoodNCM = "21069090"

	// CFOP - Código Fiscal de Operações e Prestações
	// 5102 - Venda de mercadoria adquirida ou recebida de terceiros (interna)
	DefaultCFOP = "5102"

	// CSOSN - Código de Situação da Operação do Simples Nacional
	// 102 - Tributada pelo Simples Nacional sem permissão de crédito
	// 500 - ICMS cobrado anteriormente por substituição tributária (para alimentos)
	DefaultCSOSN = "102"

	// CST - Código de Situação Tributária (para não Simples Nacional)
	// 00 - Tributada integralmente
	// 60 - ICMS cobrado anteriormente por substituição tributária
	DefaultCST = "60"

	// Origem
	// 0 - Nacional, exceto as indicadas nos códigos 3, 4, 5 e 8
	DefaultOrigem = 0

	// Unidade de medida padrão
	DefaultUnidade = "UN"

	// Alíquotas padrão (muitos estados isentam alimentação)
	DefaultAliquotaICMS = 0.0
)

// GetCSOSNForRegime returns the appropriate CSOSN/CST based on tax regime
// RegimeTributario: 1 = Simples Nacional, 2 = Simples Nacional excesso, 3 = Regime Normal
func GetCSOSNForRegime(regimeTributario int) string {
	if regimeTributario == 1 || regimeTributario == 2 {
		// Simples Nacional
		return DefaultCSOSN
	}
	// Lucro Real/Presumido
	return DefaultCST
}
