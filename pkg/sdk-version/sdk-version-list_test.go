package sdkversion_test

import (
	"encoding/json"
	"sort"

	sdkmSDKVersion "github.com/itbasis/go-tools-sdkm/pkg/sdk-version"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("SdkVersionList", func() {
	ginkgo.When("sorting", func() {
		ginkgo.It("only constructor", func() {
			sdkVersionList := sdkmSDKVersion.NewSdkVersionList(
				go1_22_5,
				go1_22_8,
				go1_22_6,
				go1_23_rc1,
				go1_23_rc2,
				go1_18_10,
				go1_3_rc1,
				go1_4_beta1,
			)
			sort.Sort(sdkVersionList)

			gomega.Expect(sdkVersionList.AsList()).To(gomega.HaveExactElements(
				// Stable
				go1_22_8,
				// Unstable
				go1_23_rc2,
				// Archived
				go1_23_rc1,
				go1_22_6,
				go1_22_5,
				go1_18_10,
				go1_4_beta1,
				go1_3_rc1,
			))
		})

		ginkgo.It("use Add()", func() {
			sdkVersionList := sdkmSDKVersion.NewSdkVersionList(
				go1_22_5,
				go1_22_6,
				go1_23_rc1,
				go1_3_rc1,
				go1_4_beta1,
			)
			sdkVersionList.Add(go1_22_8, go1_23_rc2, go1_18_10)
			sort.Sort(sdkVersionList)

			gomega.Expect(sdkVersionList.AsList()).To(gomega.HaveExactElements(
				// Stable
				go1_22_8,
				// Unstable
				go1_23_rc2,
				// Archived
				go1_23_rc1,
				go1_22_6,
				go1_22_5,
				go1_18_10,
				go1_4_beta1,
				go1_3_rc1,
			))
		})
	})

	ginkgo.When("get first", func() {
		ginkgo.It("empty list", func() {
			gomega.Expect(sdkmSDKVersion.NewSdkVersionList().First()).
				Error().To(gomega.MatchError(
				sdkmSDKVersion.ErrSDKVersionNotFound,
				gomega.ContainSubstring("list is empty"),
			))
		})

		ginkgo.It("success", func() {
			gomega.Expect(sdkmSDKVersion.NewSdkVersionList(go1_22_5, go1_22_6).First()).To(gomega.Equal(go1_22_5))
		})
	})

	ginkgo.When("MarshalJSON/UnmarshalJSON", func (){
	ginkgo.It("MarshalJSON()", func() {
		gomega.Expect(json.Marshal(sdkmSDKVersion.NewSdkVersionList(go1_22_8, go1_23_rc1))).
			To(gomega.MatchJSON(
				// language=JSON
				`
			[
			{"id": "1.22.8", "type": "stable"},
			{"id": "1.23rc1", "type": "archived"}
				]
				`,
			))
	})
	})
})
