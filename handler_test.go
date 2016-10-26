package serfer

import (
	"net"
	"testing"

	"github.com/hashicorp/serf/serf"
	log "github.com/mgutz/logxi/v1"
	"github.com/stretchr/testify/suite"
)

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestEventHandlerSuite(t *testing.T) {
	suite.Run(t, new(EventHandlerTestSuite))
}

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type EventHandlerTestSuite struct {
	suite.Suite
	Prefix  string
	Mocker  *MockEventHandler
	Handler SerfEventHandler
	Member  serf.Member
}

// Make sure that Handler and Member are set before each test
func (suite *EventHandlerTestSuite) SetupTest() {
	m := &MockEventHandler{}
	suite.Mocker = m
	suite.Prefix = "serfer"
	suite.Handler = SerfEventHandler{
		ServicePrefix:       suite.Prefix,
		UserEvent:           m,
		UnknownEventHandler: m,
		NodeJoined:          m,
		NodeLeft:            m,
		NodeFailed:          m,
		NodeReaped:          m,
		NodeUpdated:         m,
		QueryHandler:        m,
		Logger:              &log.NullLogger{},
	}

	suite.Member = serf.Member{
		Name:        "",
		Addr:        net.ParseIP("127.0.0.1"),
		Port:        9022,
		Tags:        make(map[string]string),
		Status:      serf.StatusAlive,
		ProtocolMin: serf.ProtocolVersionMin,
		ProtocolMax: serf.ProtocolVersionMax,
		ProtocolCur: serf.ProtocolVersionMax,
		DelegateMin: serf.ProtocolVersionMin,
		DelegateMax: serf.ProtocolVersionMax,
		DelegateCur: serf.ProtocolVersionMax,
	}
}

// Test NodeJoin events are processed properly
func (suite *EventHandlerTestSuite) TestNodeJoined() {

	// Create Member Event
	evt := serf.MemberEvent{
		serf.EventMemberJoin,
		[]serf.Member{suite.Member},
	}

	// Process event
	suite.Mocker.On("HandleMemberJoin", evt).Return()
	suite.Mocker.On("Reconcile", suite.Member).Return()
	suite.Handler.HandleEvent(evt)
	suite.Mocker.AssertCalled(suite.T(), "HandleMemberJoin", evt)
}

// Test NodeLeave messages are dispatched properly
func (suite *EventHandlerTestSuite) TestNodeLeave() {

	// Create Member Event
	evt := serf.MemberEvent{
		serf.EventMemberLeave,
		[]serf.Member{suite.Member},
	}

	// Process event
	suite.Mocker.On("HandleMemberLeave", evt).Return()
	suite.Mocker.On("Reconcile", suite.Member).Return()
	suite.Handler.HandleEvent(evt)
	suite.Mocker.AssertCalled(suite.T(), "HandleMemberLeave", evt)
}

// Test NodeFailed messages are dispatched properly
func (suite *EventHandlerTestSuite) TestNodeFailed() {

	// Create Member Event
	evt := serf.MemberEvent{
		serf.EventMemberFailed,
		[]serf.Member{suite.Member},
	}

	// Process event
	suite.Mocker.On("HandleMemberFailure", evt).Return()
	suite.Mocker.On("Reconcile", suite.Member).Return()
	suite.Handler.HandleEvent(evt)
	suite.Mocker.AssertCalled(suite.T(), "HandleMemberFailure", evt)
}

// Test NodeReaped messages are dispatched properly
func (suite *EventHandlerTestSuite) TestNodeReaped() {

	// Create Member Event
	evt := serf.MemberEvent{
		serf.EventMemberReap,
		[]serf.Member{suite.Member},
	}
	newMember := suite.Member
	newMember.Status = StatusReap

	// Process event
	suite.Mocker.On("HandleMemberReap", evt).Return()
	suite.Mocker.On("Reconcile", newMember).Return()
	suite.Handler.HandleEvent(evt)
	suite.Mocker.AssertCalled(suite.T(), "HandleMemberReap", evt)
}

// Test NodeUpdated messages are dispatched properly
func (suite *EventHandlerTestSuite) TestNodeUpdated() {

	// Create Member Event
	evt := serf.MemberEvent{
		serf.EventMemberUpdate,
		[]serf.Member{suite.Member},
	}

	// Process event
	suite.Mocker.On("HandleMemberUpdate", evt).Return()
	suite.Mocker.On("Reconcile", suite.Member).Return()
	suite.Handler.HandleEvent(evt)
	suite.Mocker.AssertCalled(suite.T(), "HandleMemberUpdate", evt)
}

// Test QueryEvent messages are dispatched properly
func (suite *EventHandlerTestSuite) TestQueryEvent() {

	// Create Query
	query := serf.Query{
		LTime:   serf.LamportTime(0),
		Name:    "Event",
		Payload: make([]byte, 0),
	}

	// Process event
	suite.Mocker.On("HandleQueryEvent", query).Return()
	suite.Mocker.On("Reconcile", query).Return()
	suite.Handler.HandleEvent(&query)
	suite.Mocker.AssertCalled(suite.T(), "HandleQueryEvent", query)
}

// Test nil messages are not dispatched properly
func (suite *EventHandlerTestSuite) TestNilEvent() {

	// Process event
	suite.Handler.HandleEvent(nil)
	suite.Mocker.AssertNotCalled(suite.T(), "HandleMemberJoin")
	suite.Mocker.AssertNotCalled(suite.T(), "HandleMemberLeave")
	suite.Mocker.AssertNotCalled(suite.T(), "HandleMemberFailure")
	suite.Mocker.AssertNotCalled(suite.T(), "HandleMemberUpdate")
	suite.Mocker.AssertNotCalled(suite.T(), "HandleMemberReap")
	suite.Mocker.AssertNotCalled(suite.T(), "HandleUknownEvent")
	suite.Mocker.AssertNotCalled(suite.T(), "HandleUserEvent")
	suite.Mocker.AssertNotCalled(suite.T(), "HandleQueryEvent")
}

// Test unknown messages are not dispatched properly
func (suite *EventHandlerTestSuite) TestUnknownEvent() {

	// Process event
	t1 := &MockEvent{Name: "UnknownType", Type: serf.EventType(-1)}
	t1.On("EventType").Return()
	suite.Handler.HandleEvent(t1)

	// Test Assertions
	t1.AssertCalled(suite.T(), "EventType")
	suite.Mocker.AssertNotCalled(suite.T(), "HandleMemberJoin")
	suite.Mocker.AssertNotCalled(suite.T(), "HandleMemberLeave")
	suite.Mocker.AssertNotCalled(suite.T(), "HandleMemberFailure")
	suite.Mocker.AssertNotCalled(suite.T(), "HandleMemberUpdate")
	suite.Mocker.AssertNotCalled(suite.T(), "HandleMemberReap")
	suite.Mocker.AssertNotCalled(suite.T(), "HandleUknownEvent")
	suite.Mocker.AssertNotCalled(suite.T(), "HandleUserEvent")
	suite.Mocker.AssertNotCalled(suite.T(), "HandleQueryEvent")
}

// Test unknown user events are dispatched properly
func (suite *EventHandlerTestSuite) TestUserEvent_UnknownEvent() {

	// Create Member Event
	evt := serf.UserEvent{
		LTime:    serf.LamportTime(0),
		Name:     "unk",
		Payload:  make([]byte, 0),
		Coalesce: false,
	}

	// Process event
	suite.Mocker.On("HandleUnknownEvent", evt).Return()
	suite.Handler.HandleEvent(evt)
	suite.Mocker.AssertCalled(suite.T(), "HandleUnknownEvent", evt)
}

// Test user events are dispatched properly
func (suite *EventHandlerTestSuite) TestUserEvent() {

	// Create Member Event
	evt := serf.UserEvent{
		LTime:    serf.LamportTime(0),
		Name:     suite.Prefix + ":user-event",
		Payload:  make([]byte, 0),
		Coalesce: false,
	}
	modified := evt
	modified.Name = "user-event"

	// Process event
	suite.Mocker.On("HandleUserEvent", modified).Return()
	suite.Handler.HandleEvent(evt)
	suite.Mocker.AssertCalled(suite.T(), "HandleUserEvent", modified)
}
