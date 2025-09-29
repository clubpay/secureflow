package v5_4

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"time"

	"github.com/clubpay-pos-worker/sdk-go/v2"
	"github.com/clubpay-pos-worker/sdk-go/v2/qlub"
	"github.com/clubpay/qlubkit-go/telemetry/log"
	qtrace "github.com/clubpay/qlubkit-go/telemetry/trace"
	"github.com/hooklift/gowsdl/soap"
	"golang.org/x/sync/singleflight"

	"github.com/clubpay-pos-worker/oraclets/micros"
	"github.com/clubpay-pos-worker/oraclets/micros/v5_4/api"
)

//go:generate rm -rf ./api
//go:generate mkdir -p ./api
//go:generate go run github.com/hooklift/gowsdl/cmd/gowsdl -o protocol.go -p api -i ResPosApiWeb.wsdl

type acl struct {
	employeeIdNumber string
	base             api.ResPosApiWebServiceSoap

	cache    micros.Cache
	sfTables singleflight.Group

	lg sdk.Logger
}

const urlRoot = "ResPosApiWeb/ResPosApiWeb.asmx"

func NewPos(baseUrl url.URL, opts ...micros.Option) micros.Pos {
	m := &acl{
		cache: micros.NoCache,
		base:  api.NewResPosApiWebServiceSoap(soap.NewClient(fmt.Sprintf("%s/%s", baseUrl.String(), urlRoot))),

		lg: log.DefaultLogger,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *acl) SetLogger(lg sdk.Logger) {
	m.lg = lg
}

func (m *acl) SetStore(_ sdk.Store) {}

func (m *acl) SetAPI(_ sdk.API) {
}

func (m *acl) Authenticate(creds string) {
	m.employeeIdNumber = creds
}

func (m *acl) SetCache(cache micros.Cache) {
	m.cache = cache
}

func (m *acl) AutofillCache(interval time.Duration) {
	go m.autofillCache(interval)
}

func (m *acl) FetchTable(ctx context.Context, revenueCenter, floor, section, tableID string) (*micros.Table, error) {
	m.lg.InfoCtx(
		ctx, "enter FetchTable method",
		log.String(qtrace.QlubTableRevenueCenter, revenueCenter),
		log.String("table_floor", floor),
		log.String(qtrace.QlubTableSection, section),
		log.String(qtrace.QlubTableID, tableID),
	)

	if tbl, _ := m.cache.GetTable(ctx, micros.TableUnique(revenueCenter, floor, section, tableID)); tbl != nil {
		m.lg.DebugCtx(ctx, "found in cache")

		return tbl, nil
	}

	cc, err := m.fetchAndCacheAllChecks(ctx)
	if err != nil {
		m.lg.ErrorCtx(ctx, "error occurred while fetching all tables", log.Error(err))

		return nil, err
	}

	for _, check := range cc {
		m.lg.DebugCtx(ctx, "processing check", log.Reflect("check", check))
		if check == nil {
			continue
		}

		if tbl := checkSummaryWithSeatsToTablePtr(check); tbl != nil && tbl.Match(revenueCenter, floor, section, tableID) {
			m.lg.DebugCtx(
				ctx, "found table in result set",
				log.Int32("table object number", check.CheckTableObjectNum),
			)

			return tbl, nil
		}
	}

	m.lg.DebugCtx(ctx, "no check table matched")

	return nil, micros.ErrNotFound
}

/*
func (m *acl) FetchOrder(ctx context.Context, orderID string) (*micros.Order, error) {
	m.lg.InfoCtx(ctx, "enter FetchOrder method", log.String(qtrace.QlubOrderID, orderID))

	checkSeq, errParse := strconv.Atoi(orderID)
	if errParse != nil {
		m.lg.ErrorCtx(ctx, "illegal orderId", log.Error(errParse))

		return nil, errParse
	}

	m.lg.DebugCtx(ctx, "calling getCheckSummary", log.Int("check sequence", checkSeq))
	checkSummary, err := m.getCheckSummary(ctx, int32(checkSeq))
	if err != nil {
		m.lg.ErrorCtx(ctx, "error on getCheckSummary", log.Error(err))

		return nil, err
	}

	//m.lg.DebugCtx(ctx, "calling getPrintedCheck: checkSeq(%d)", checkSeq)
	//checkPrinted, err := m.getPrintedCheck(ctx, int32(checkSeq))
	//if err != nil {
	//	span.RecordError(err)
	//	m.lg.ErrorCtx(ctx, "error on getPrintedCheck: %s", err)
	//
	//	return nil, err
	//}

	o := checkSummaryWithSeatsToOrderPtr(
		checkSummary,
		//guestCheckToOrderModifierPtrVal(checkPrinted.PGuestCheck),
		//totalsResponseToOrderModifierPtrVal(checkPrinted.PTotalsResponse),
	)

	return o, nil
}
*/

func (m *acl) FetchOrder(ctx context.Context, orderID, tenderMediaNumber, revenueBuckets, qlubDiscountCodes string) (*qlub.Order, error) {
	panic("not implemented")
}

func (m *acl) ApplyPayment(ctx context.Context, orderID, txnID, total, bill, tip string) error {
	m.lg.InfoCtx(
		ctx, "enter ApplyPayment method",
		log.String(qtrace.QlubOrderID, orderID),
		log.String(qtrace.QlubPaymentTxnID, txnID),
		log.String(qtrace.QlubPaymentAmountBill, bill),
		log.String(qtrace.QlubPaymentAmountTip, tip),
	)

	checkSeq, errParse := strconv.Atoi(orderID)
	if errParse != nil {
		m.lg.ErrorCtx(ctx, "illegal orderId", log.Error(errParse))

		return errParse
	}

	if checkSeq < math.MinInt32 || checkSeq > math.MaxInt32 {
		m.lg.ErrorCtx(ctx, "orderId out of int32 range", log.Int("checkSeq", checkSeq))
		return fmt.Errorf("orderId %d out of int32 range", checkSeq)
	}

	checkSeqCasted := int32(checkSeq) // nosemgrep
	err := m.addPayment(ctx, checkSeqCasted,txnID, bill, tip) // nosemgrep
	if err != nil {
		m.lg.ErrorCtx(ctx, "error on addPayment", log.Error(err))

		return err
	}

	return nil
}

func (m *acl) autofillCache(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {
		// call auto-fillers
		go func() { _, _ = m.fetchAndCacheAllChecks(context.TODO()) }()

		<-ticker.C
	}
}

func (m *acl) fetchAndCacheAllChecks(ctx context.Context) ([]*api.ResPosAPI_CheckSummaryWithSeats, error) {
	m.lg.InfoCtx(ctx, "enter fetchAndCacheAllChecks method")

	v, err, _ := m.sfTables.Do("allChecks", func() (interface{}, error) {
		req := new(api.GetOpenChecksWithSeats)

		m.lg.DebugCtx(ctx, "GetOpenChecksWithSeatsContext request ready", log.Reflect("request", req))
		rsp, err := m.base.GetOpenChecksWithSeatsContext(ctx, req)
		if err != nil {
			m.lg.ErrorCtx(ctx, "error on GetOpenChecksWithSeatsContext", log.Error(err))

			return nil, err
		}

		m.lg.DebugCtx(ctx, "GetOpenChecksWithSeatsContext response", log.Reflect("response", rsp))
		if rsp.PpCheckSummaryArrayWithSeats == nil {
			m.lg.ErrorCtx(ctx, "empty response")

			return nil, micros.ErrNotFound
		}

		cc := rsp.PpCheckSummaryArrayWithSeats.ResPosAPI_CheckSummaryWithSeats

		// fill cache
		if err := m.cache.AddTables(ctx, checkSummaryWithSeatsToTableListPtrVal(cc)...); err != nil {
			m.lg.ErrorCtx(ctx, "error occurred while tried to caching tables", log.Error(err))
		}

		return cc, nil
	})

	m.lg.DebugCtx(ctx, "fetch all checks single-flight returned")
	if err != nil {
		m.lg.ErrorCtx(ctx, "error on fetching all checks", log.Error(err))

		return nil, err
	}

	cc, _ := v.([]*api.ResPosAPI_CheckSummaryWithSeats)

	return cc, err
}

func (m *acl) getCheckSummary(ctx context.Context, checkSeq int32) (*api.ResPosAPI_CheckSummaryWithSeats, error) {
	m.lg.DebugCtx(ctx, "enter getCheckSummary method", log.Int32("check sequence", checkSeq))

	req := new(api.GetCheckSummaryWithSeats)
	req.EmployeeIdNum = m.employeeIdNumber
	req.CheckSeq = checkSeq

	m.lg.DebugCtx(ctx, "GetCheckSummaryWithSeatsContext request ready", log.Reflect("request", req))
	rsp, err := m.base.GetCheckSummaryWithSeatsContext(ctx, req)
	if err != nil {
		m.lg.ErrorCtx(ctx, "error on GetCheckSummaryWithSeatsContext: %s", log.Error(err))

		return nil, err
	}

	m.lg.DebugCtx(ctx, "GetCheckSummaryWithSeatsContext response", log.Reflect("response", rsp))
	if rsp == nil {
		return nil, micros.ErrNotFound
	} else if rsp.PCheckSummary == nil {
		return nil, micros.ErrNotFound
	}

	return rsp.PCheckSummary, nil
}

func (m *acl) getPrintedCheck(ctx context.Context, checkSeq int32) (*api.GetPrintedCheckResponse, error) {
	m.lg.DebugCtx(ctx, "enter getPrintedCheck method", log.Int32("check sequence", checkSeq))

	req := new(api.GetPrintedCheck)
	req.CheckSeq = checkSeq

	m.lg.DebugCtx(ctx, "GetPrintedCheckContext request ready", log.Reflect("request", req))
	rsp, err := m.base.GetPrintedCheckContext(ctx, req)
	if err != nil {
		m.lg.ErrorCtx(ctx, "error on GetPrintedCheckContext", log.Error(err))

		return nil, err
	}

	m.lg.DebugCtx(ctx, "GetPrintedCheckContext response", log.Reflect("response", rsp))
	if rsp == nil {
		return nil, micros.ErrNotFound
	} else if rsp.PGuestCheck == nil {
		return nil, micros.ErrNotFound
	}

	return rsp, nil
}

func (m *acl) addPayment(ctx context.Context, checkSeq int32, txnID, bill, tip string) error {
	req := new(api.AddToExistingCheckEx)
	req.PGuestCheck = &api.ResPosAPI_GuestCheck{
		CheckSeq: checkSeq,
	}
	req.PTmedDetail = &api.ResPosAPI_TmedDetailItemEx{
		TmedObjectNum:      0,
		TmedPartialPayment: "",
		TmedReference:      "",
		TmedEPayment: &api.ResPosAPI_EPayment{
			PaymentCommand:      nil,
			AccountDataSource:   nil,
			AccountType:         nil,
			AcctNumber:          "",
			ExpirationDate:      soap.XSDDateTime{},
			AuthorizationCode:   "",
			StartDate:           soap.XSDDateTime{},
			IssueNumber:         0,
			Track1Data:          "",
			Track2Data:          "",
			Track3Data:          "",
			BaseAmount:          "",
			TipAmount:           "",
			CashBackAmount:      "",
			KeySerialNum:        "",
			DeviceId:            "",
			PinBlock:            "",
			CVVNumber:           "",
			AddressVerification: "",
			InterfaceName:       "",
			SvcResponse:         "",
			SvcAccountType:      "",
		},
	}

	_, err := m.base.AddToExistingCheckExContext(ctx, req)
	if err != nil {
		m.lg.ErrorCtx(ctx, "error on AddToExistingCheckExContext", log.Error(err))

		return err
	}

	return err
}
