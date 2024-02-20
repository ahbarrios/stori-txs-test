package template

import (
	_ "embed"
	"html/template"
	"strings"

	"github.com/ahbarrios/stori-txn-test/internal/tx/agg"
)

//go:embed summary.html.tpl
var summarySrcTpl string

var summaryTpl *template.Template

// mso it will be used to add the Microsoft Outlook conditional comment to the [template.HTML] template safely
// This code helps rendering on Windows versions of Outlook desktop.
var mso = template.HTML(`<!--[if mso]> 
    <noscript> 
    <xml> 
    <o:OfficeDocumentSettings> 
    <o:PixelsPerInch>96</o:PixelsPerInch> 
    </o:OfficeDocumentSettings> 
    </xml> 
    </noscript> 
    <![endif]-->`)

// NewSummaryBody it will create a new [email.SummaryEmailBodyFunc] compliant function that will be used to create the summary email using the embed [HTML template].
// The above link offers an example of the expected [HTML template] that will be used to create the email body.
//
// [HTML template]: https://github/ahbarrios/stori-txn-test/internal/email/template/summary_example_test.html
func NewSummaryBody(balance agg.Balance, credit agg.AvgCredit, debit agg.AvgDebit, countByMonth agg.MonthlySummary) (*strings.Reader, error) {
	var body strings.Builder
	if err := summaryTpl.Execute(&body, struct {
		Mso          template.HTML
		Balance      agg.Balance
		AvgCredit    agg.AvgCredit
		AvgDebit     agg.AvgDebit
		CountByMonth agg.MonthlySummary
	}{
		Mso:          mso,
		Balance:      balance,
		AvgCredit:    credit,
		AvgDebit:     debit,
		CountByMonth: countByMonth,
	}); err != nil {
		return nil, err
	}
	return strings.NewReader(body.String()), nil
}

func init() {
	summaryTpl = template.Must(template.New("summary").Parse(summarySrcTpl))
}
