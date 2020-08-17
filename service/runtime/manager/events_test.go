package manager

import (
	"testing"
	"time"

	"github.com/micro/go-micro/v3/runtime"
	"github.com/micro/micro/v3/internal/namespace"
	"github.com/micro/micro/v3/profile"
	muruntime "github.com/micro/micro/v3/service/runtime"
)

func TestEvents(t *testing.T) {
	// an event is passed through this channel from the test runtime everytime a method is called,
	// this is done since events ae processed async
	eventChan := make(chan *runtime.Service)

	profile.Test.Setup(nil)
	rt := &testRuntime{events: eventChan}
	muruntime.DefaultRuntime = rt
	m := New().(*manager)

	// set the eventPollFrequency to 10ms so events are processed immediately
	eventPollFrequency = time.Millisecond * 10
	go m.watchEvents()

	// timeout async tests after 500ms
	timeout := time.Millisecond * 500

	// the service that should be passed to the runtime
	testSrv := &runtime.Service{Name: "foo", Version: "latest"}
	opts := &runtime.CreateOptions{Namespace: namespace.DefaultNamespace}

	t.Run("Create", func(t *testing.T) {
		defer rt.Reset()

		if err := m.publishEvent(runtime.CreatedEvent, testSrv, opts); err != nil {
			t.Errorf("Unexpected error when publishing events: %v", err)
		}

		select {
		case srv := <-eventChan:
			if srv.Name != testSrv.Name || srv.Version != testSrv.Version {
				t.Errorf("Incorrect service passed to the runtime")
			}
		case <-time.After(timeout):
			t.Fatalf("The runtime wasn't called")
		}

		if rt.createCount != 1 {
			t.Errorf("Expected runtime create to be called 1 time but was actually called %v times", rt.createCount)
		}
	})

	t.Run("Update", func(t *testing.T) {
		defer rt.Reset()

		if err := m.publishEvent(runtime.UpdatedEvent, testSrv, opts); err != nil {
			t.Errorf("Unexpected error when publishing events: %v", err)
		}

		select {
		case srv := <-eventChan:
			if srv.Name != testSrv.Name || srv.Version != testSrv.Version {
				t.Errorf("Incorrect service passed to the runtime")
			}
		case <-time.After(timeout):
			t.Fatalf("The runtime wasn't called")
		}

		if rt.updateCount != 1 {
			t.Errorf("Expected runtime update to be called 1 time but was actually called %v times", rt.updateCount)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		defer rt.Reset()

		if err := m.publishEvent(runtime.DeletedEvent, testSrv, opts); err != nil {
			t.Errorf("Unexpected error when publishing events: %v", err)
		}

		select {
		case srv := <-eventChan:
			if srv.Name != testSrv.Name || srv.Version != testSrv.Version {
				t.Errorf("Incorrect service passed to the runtime")
			}
		case <-time.After(timeout):
			t.Fatalf("The runtime wasn't called")
		}

		if rt.deleteCount != 1 {
			t.Errorf("Expected runtime delete to be called 1 time but was actually called %v times", rt.deleteCount)
		}
	})
}
