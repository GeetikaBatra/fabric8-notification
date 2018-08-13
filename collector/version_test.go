package collector_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fabric8-services/fabric8-notification/auth"
	authApi "github.com/fabric8-services/fabric8-notification/auth/api"
	"github.com/fabric8-services/fabric8-notification/collector"
	"github.com/fabric8-services/fabric8-notification/testsupport"
	"github.com/fabric8-services/fabric8-notification/wit"
	witApi "github.com/fabric8-services/fabric8-notification/wit/api"
	"github.com/stretchr/testify/assert"
)

func TestCVEResolver(t *testing.T) {
	witServer := createServer(serveWITRequest)
	authServer := createServer(serveAuthRequest)

	witURL := "http://" + witServer.Listener.Addr().String() + "/"
	authURL := "http://" + authServer.Listener.Addr().String() + "/"

	witClient, authClient := createLocalClient(t, witURL, authURL)

	cveResolver := collector.NewCVEResolver(authClient, witClient)
	codebaseURL := "git@github.com:testrepo/testproject1.git"
	recvs, _, err := cveResolver(context.Background(), codebaseURL)

	assert.Nil(t, err)
	assert.NotNil(t, recvs)
	assert.Equal(t, 2, len(recvs))
	checkEmails(t, recvs, "testuser1@redhat.com", "testuser2@redhat.com")
}