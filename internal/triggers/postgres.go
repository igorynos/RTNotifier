package triggers

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

type QueryTrigger struct {
	NameValue, Query string
	DB               *pgxpool.Pool
	Format           func([]string) string
}

func (t QueryTrigger) Name() string { return t.NameValue }
func (t QueryTrigger) Check(ctx context.Context) (bool, string, error) {
	rows, err := t.DB.Query(ctx, t.Query)
	if err != nil {
		return false, "", err
	}
	defer rows.Close()
	fields := rows.FieldDescriptions()
	var lines []string
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return false, "", err
		}
		parts := make([]string, len(values))
		for i, v := range values {
			parts[i] = fmt.Sprintf("%s=%v", fields[i].Name, v)
		}
		lines = append(lines, strings.Join(parts, " "))
	}
	if err = rows.Err(); err != nil {
		return false, "", err
	}
	if len(lines) == 0 {
		return false, "", nil
	}
	if t.Format != nil {
		return true, t.Format(lines), nil
	}
	return true, "🚨 " + t.NameValue + "\n" + strings.Join(lines, "\n"), nil
}
func All(db *pgxpool.Pool) []QueryTrigger {
	return []QueryTrigger{{"provider-issues", providerIssues, db, nil}, {"client-issues", clientIssues, db, nil}, {"traffic-stopped", trafficStopped, db, nil}, {"client-launched", clientLaunched, db, nil}, {"low-conversion", lowConversion, db, nil}}
}

const providerIssues = `SELECT pg.original_merchant_id provider_id,pp.provider_code,COUNT(*) FILTER(WHERE c.status='process_pending') pending,COUNT(*) FILTER(WHERE c.status='process_failed') failed FROM pc_com_payment_invoices c JOIN pc_pg_payments pg ON c.active_payment_id=pg.id JOIN pay_invoices pi ON pg.external_id=pi.id JOIN pay_requests pr ON pi.active_request_id=pr.id JOIN pay_payments pp ON pr.active_payment_id=pp.id WHERE c.created>=NOW()-INTERVAL '2 hours' GROUP BY pg.original_merchant_id,pp.provider_code HAVING COUNT(*) FILTER(WHERE c.status='process_pending')>=15 OR COUNT(*) FILTER(WHERE c.status='process_failed')>=15`
const clientIssues = `SELECT a.id account_id,a.name,COUNT(*) FILTER(WHERE i.status='process_pending') pending,COUNT(*) FILTER(WHERE i.status='process_failed') failed FROM pc_com_accounts a JOIN pc_com_payment_invoices i ON i.commerce_account_id=a.id WHERE i.created>=NOW()-INTERVAL '2 hours' GROUP BY a.id,a.name HAVING COUNT(*) FILTER(WHERE i.status='process_pending')>=10 OR COUNT(*) FILTER(WHERE i.status='process_failed')>=10`
const trafficStopped = `WITH recent AS(SELECT commerce_account_id,created FROM pc_com_payment_invoices UNION ALL SELECT commerce_account_id,created FROM pc_com_payout_invoices) SELECT a.id,a.name FROM pc_com_accounts a JOIN recent r ON r.commerce_account_id=a.id GROUP BY a.id,a.name HAVING COUNT(*) FILTER(WHERE r.created>=NOW()-INTERVAL '24 hours')>=10 AND COUNT(*) FILTER(WHERE r.created>=NOW()-INTERVAL '3 hours')=0`
const clientLaunched = `SELECT a.id,a.name,COUNT(DISTINCT pi.id)+COUNT(DISTINCT po.id) invoice_count FROM pc_com_accounts a LEFT JOIN pc_com_payment_invoices pi ON a.id=pi.commerce_account_id LEFT JOIN pc_com_payout_invoices po ON a.id=po.commerce_account_id WHERE a.created>=CURRENT_DATE-INTERVAL '14 days' GROUP BY a.id HAVING COUNT(*) FILTER(WHERE pi.status='processed')+COUNT(*) FILTER(WHERE po.status='processed')>=10`
const lowConversion = `SELECT pg.original_merchant_id,pp.provider_code,ROUND(100.0*COUNT(*) FILTER(WHERE c.status='processed')/NULLIF(COUNT(*),0),2) conversion FROM pc_com_payment_invoices c JOIN pc_pg_payments pg ON c.active_payment_id=pg.id JOIN pay_invoices pi ON pg.external_id=pi.id JOIN pay_requests pr ON pi.active_request_id=pr.id JOIN pay_payments pp ON pr.active_payment_id=pp.id WHERE c.created>=NOW()-INTERVAL '24 hours' GROUP BY pg.original_merchant_id,pp.provider_code HAVING 100.0*COUNT(*) FILTER(WHERE c.status='processed')/NULLIF(COUNT(*),0)<50`
