package pos

const style = `
<style>
    body {
        font-family: 'Courier New', Courier, monospace;
        font-size: 12px;
        width: 80mm;
        margin: 0;
        padding: 5px;
    }
    .header, .footer {
        text-align: center;
        margin-bottom: 10px;
    }
    .bold {
        font-weight: bold;
    }
    .align-left {
        text-align: left;
    }
    .align-right {
        text-align: right;
    }
    .row {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
    }
    .col-name {
        flex-grow: 1;
        text-align: left;
        padding-right: 10px;
    }
    .col-price {
        white-space: nowrap;
        text-align: right;
    }
    .divider {
        border-top: 1px dashed #000;
        margin: 5px 0;
    }
    .item {
        margin-bottom: 5px;
    }
    .obs {
        font-style: italic;
        margin-left: 10px;
        font-size: 0.9em;
    }
</style>
`

const KitchenTicketTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    ` + style + `
</head>
<body>
    <div class="header bold">
        {{if .Company}}{{.Company.TradeName}}{{end}}
    </div>
    <div class="header">Cozinha ({{.Quantity}} itens)</div>
    <div class="divider"></div>
    
    <div class="item bold">
        {{if .Category}}
            {{.Category.Name}}
        {{end}}
        {{if .Size}}
            - {{.Size}}
        {{end}}
    </div>
    <div class="divider"></div>

    {{range .Items}}
    <div class="item">
        <div class="row bold">
            <span>{{.Quantity}}x {{.Name}}</span>
        </div>
        {{range .AdditionalItems}}
        <div class="row">
            <span>+ {{.Quantity}}x {{.Name}}</span>
        </div>
        {{end}}
        {{range .RemovedItems}}
        <div class="row">
            <span>- {{.}}</span>
        </div>
        {{end}}
        {{if .Observation}}
        <div class="obs">
            Obs: {{.Observation}}
        </div>
        {{end}}
    </div>
    {{end}}

    {{if .ComplementItem}}
    <div class="divider"></div>
    <div class="item">
        <div class="bold">Complemento:</div>
        <div>{{.ComplementItem.Name}}</div>
    </div>
    {{end}}

    {{if .Observation}}
    <div class="divider"></div>
    <div class="obs bold">
        Obs Geral: {{.Observation}}
    </div>
    {{end}}

    <div class="footer">
        {{if .StartAt}}
            Agendado: {{.StartAt.Format "15:04"}}
        {{end}}
    </div>
</body>
</html>
`

const OrderReceiptTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    ` + style + `
</head>
<body>
    <div class="header bold">
        {{if .Company}}{{.Company.TradeName}}{{end}}
    </div>
    {{if .Company}}
    <div class="header">
        {{if .Company.Cnpj}}CNPJ: {{.Company.Cnpj}}<br>{{end}}
        {{if .Company.Address}}
            {{.Company.Address.Street}}, {{.Company.Address.Number}} - {{.Company.Address.Neighborhood}}<br>
        {{end}}
        {{if .Company.Contacts}}
            Contato: {{index .Company.Contacts 0}}
        {{end}}
    </div>
    {{end}}
    <div class="divider"></div>
    <div class="header bold">
        PEDIDO #{{.OrderNumber}}
    </div>
    <div class="header">
        {{if .PendingAt}}
            {{.PendingAt.Format "02/01/2006 15:04"}}
        {{end}}
    </div>
    
    <div class="divider"></div>

    {{if .Delivery}}
        <div class="header bold">ENTREGA</div>
        {{if .Delivery.Client}}
            <div>Cliente: {{.Delivery.Client.Name}}</div>
        {{end}}
        {{if .Delivery.Address}}
            <div>{{.Delivery.Address.Street}}, {{.Delivery.Address.Number}}</div>
            <div>{{.Delivery.Address.Neighborhood}}</div>
        {{end}}
        <div class="divider"></div>
    {{end}}

    {{if .Pickup}}
        <div class="header bold">RETIRADA</div>
        {{if .Pickup.Name}}
            <div>Cliente: {{.Pickup.Name}}</div>
        {{end}}
        <div class="divider"></div>
    {{end}}

    {{if .Table}}
        <div class="header bold">MESA</div>
        {{if .Table.Name}}
            <div>Mesa: {{.Table.Name}}</div>
        {{end}}
        <div class="divider"></div>
    {{end}}

    {{if .Attendant}}
        {{if .Attendant.User}}
        <div class="header">
            Atendente: {{.Attendant.User.Name}}
        </div>
        <div class="divider"></div>
        {{end}}
    {{end}}

    {{if .Observation}}
        <div class="item bold">
            Obs: {{.Observation}}
        </div>
        <div class="divider"></div>
    {{end}}

    <div class="items">
        {{range .GroupItems}}
            <div class="item bold">
                {{if .Category}}{{.Category.Name}}{{end}} {{if .Size}}- {{.Size}}{{end}} 	Qtd: {{.Quantity}}
            </div>
            {{range .Items}}
                <div class="row">
                    <span class="col-name">{{.Quantity}}x {{.Name}}</span>
                    <span class="col-price">{{ multiply .Price .Quantity | formatMoney }}</span>
                </div>
                {{range .AdditionalItems}}
                    <div class="row">
                        <span class="col-name">+ {{.Quantity}}x {{.Name}}</span>
                        <span class="col-price">{{formatMoney .TotalPrice}}</span>
                    </div>
                {{end}}
                {{range .RemovedItems}}
                    <div class="row">
                        <span class="col-name">- {{.}}</span>
                    </div>
                {{end}}
                {{if .Observation}}
                    <div class="obs">Obs: {{.Observation}}</div>
                {{end}}
            {{end}}
            
            {{if .ComplementItem}}
                <div class="row">
                    <span class="col-name">Comp: {{.ComplementItem.Name}}</span>
                    <span class="col-price">{{formatMoney .ComplementItem.TotalPrice}}</span>
                </div>
            {{end}}

            <div class="row bold">
                <span class="col-name">Subtotal:</span>
                <span class="col-price">{{formatMoney .TotalPrice}}</span>
            </div>

            <div class="divider"></div>
        {{end}}
    </div>

    <div class="totals bold">
        <div class="row">
            <span class="col-name">SUBTOTAL:</span>
            <span class="col-price">{{formatMoney .SubTotal}}</span>
        </div>
        <div class="row">
            <span class="col-name">TOTAL:</span>
            <span class="col-price">{{formatMoney .Total}}</span>
        </div>
    </div>

    <div class="payments">
        {{range .Payments}}
            <div class="row">
                <span class="col-name">{{.Method}}</span>
                <span class="col-price">{{formatMoney .TotalPaid}}</span>
            </div>
        {{end}}
    </div>
    
    {{if .TotalChange}}
    <div class="row">
        <span class="col-name">Troco:</span>
        <span class="col-price">{{formatMoney .TotalChange}}</span>
    </div>
    {{end}}

    <div class="footer">
        <p>Obrigado pela preferência!</p>
    </div>
</body>
</html>
`

const ShiftReportTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    ` + style + `
</head>
<body>
    <div class="header bold">
        RELATÓRIO DE TURNO
    </div>
    <div class="header">
        Abertura: {{if .OpenedAt}}{{.OpenedAt.Format "02/01/2006 15:04"}}{{end}}<br>
        {{if .ClosedAt}}Fechamento: {{.ClosedAt.Format "02/01/2006 15:04"}}{{end}}
    </div>
    <div class="divider"></div>

    <div class="header bold">RESUMO DE VENDAS</div>
    <div class="row">
        <span>Pedidos Finalizados:</span>
        <span>{{.TotalOrdersFinished}}</span>
    </div>
    <div class="row">
        <span>Pedidos Cancelados:</span>
        <span>{{.TotalOrdersCancelled}}</span>
    </div>
    <div class="row bold">
        <span>TOTAL GERAL:</span>
        <span>{{formatMoney .TotalSales}}</span>
    </div>

    <div class="divider"></div>
    <div class="header bold">VENDAS POR CATEGORIA</div>
    {{range $cat, $val := .SalesByCategory}}
    <div class="row">
        <span class="col-name">{{$cat}}</span>
        <span class="col-price">{{formatMoney $val}}</span>
    </div>
    {{end}}

    <div class="divider"></div>
    <div class="header bold">PRODUTOS VENDIDOS</div>
    {{range $cat, $qty := .ProductsSoldByCategory}}
    <div class="row">
        <span class="col-name">{{$cat}}</span>
        <span class="col-price">{{$qty}}</span>
    </div>
    {{end}}

    <div class="divider"></div>
    <div class="header bold">PAGAMENTOS</div>
    {{range .Payments}}
    <div class="row">
        <span class="col-name">{{.Method}}</span>
        <span class="col-price">{{formatMoney .TotalPaid}}</span>
    </div>
    {{end}}

    <div class="divider"></div>
    <div class="header bold">LISTA DE PEDIDOS</div>
    {{range .Orders}}
    <div class="row">
        <span class="col-name">Pedido #{{.OrderNumber}}</span>
        <span class="col-price">{{formatMoney .SubTotal}}</span>
        <span class="col-price">{{formatMoney .Total}}</span>
    </div>
    {{end}}

    <div class="divider"></div>
    <div class="footer">
        <p>Relatório gerado em: {{now}}</p>
    </div>
</body>
</html>
`
