package filestorage

import (
	"context"
	"path"
	"time"

	sdkmPlugin "github.com/itbasis/go-tools-sdkm/pkg/plugin"
	"github.com/jonboulle/clockwork"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.DescribeTable("validate", func(date string, want bool) {
	var tTest, errTest = time.Parse(UpdatedFormat, "2025-01-01 00:00:01")
	gomega.Expect(errTest).To(gomega.Succeed())

	var t, err = time.Parse(UpdatedFormat, date)
	gomega.Expect(err).To(gomega.Succeed())

	var (
		ctx     = clockwork.AddToContext(context.Background(), clockwork.NewFakeClockAt(t))
		storage = NewFileCacheStorage(path.Join("testdata", "001"), sdkmPlugin.ID("go")).(*_fileStorage)
		m       = model{
			Updated: updated(tTest),
		}
	)

	gomega.Expect(storage.validate(ctx, m)).To(gomega.Equal(want))
},
	ginkgo.Entry(nil, "2024-12-31 00:00:00", true),
	ginkgo.Entry(nil, "2024-12-31 00:00:01", true),
	ginkgo.Entry(nil, "2025-01-01 00:00:00", true),
	ginkgo.Entry(nil, "2025-01-01 00:00:01", true),
	ginkgo.Entry(nil, "2025-01-01 00:00:02", true),
	ginkgo.Entry(nil, "2025-01-02 00:00:01", true),
	ginkgo.Entry(nil, "2025-01-02 00:00:02", false),
)
