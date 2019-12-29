package tenable

import "encoding/xml"

// ScansExportNessusData export-scans will output a XML like this
type ScansExportNessusData struct {
	XMLName xml.Name `xml:"NessusClientData_v2"`
	Text    string   `xml:",chardata"`
	Policy  struct {
		Text           string `xml:",chardata"`
		PolicyName     string `xml:"policyName"`
		PolicyComments string `xml:"policyComments"`
		Preferences    struct {
			Text              string `xml:",chardata"`
			ServerPreferences struct {
				Text       string `xml:",chardata"`
				Preference []struct {
					Text  string `xml:",chardata"`
					Name  string `xml:"name"`
					Value string `xml:"value"`
				} `xml:"preference"`
			} `xml:"ServerPreferences"`
			PluginsPreferences struct {
				Text string `xml:",chardata"`
				Item struct {
					Text             string `xml:",chardata"`
					PluginName       string `xml:"pluginName"`
					PluginId         string `xml:"pluginId"`
					FullName         string `xml:"fullName"`
					PreferenceName   string `xml:"preferenceName"`
					PreferenceType   string `xml:"preferenceType"`
					PreferenceValues string `xml:"preferenceValues"`
					SelectedValue    string `xml:"selectedValue"`
				} `xml:"item"`
			} `xml:"PluginsPreferences"`
		} `xml:"Preferences"`
		FamilySelection struct {
			Text       string `xml:",chardata"`
			FamilyItem []struct {
				Text       string `xml:",chardata"`
				FamilyName string `xml:"FamilyName"`
				Status     string `xml:"Status"`
			} `xml:"FamilyItem"`
		} `xml:"FamilySelection"`
		IndividualPluginSelection struct {
			Text       string `xml:",chardata"`
			PluginItem struct {
				Text       string `xml:",chardata"`
				PluginId   string `xml:"PluginId"`
				PluginName string `xml:"PluginName"`
				Family     string `xml:"Family"`
				Status     string `xml:"Status"`
			} `xml:"PluginItem"`
		} `xml:"IndividualPluginSelection"`
	} `xml:"Policy"`
	Report struct {
		Text       string `xml:",chardata"`
		Name       string `xml:"name,attr"`
		Cm         string `xml:"cm,attr"`
		ReportHost []struct {
			Text           string `xml:",chardata"`
			Name           string `xml:"name,attr"`
			HostProperties struct {
				Text string `xml:",chardata"`
				Tag  []struct {
					Text string `xml:",chardata"`
					Name string `xml:"name,attr"`
				} `xml:"tag"`
			} `xml:"HostProperties"`
			ReportItem []struct {
				Text                   string `xml:",chardata"`
				Severity               string `xml:"severity,attr"`
				Port                   string `xml:"port,attr"`
				PluginFamily           string `xml:"pluginFamily,attr"`
				PluginName             string `xml:"pluginName,attr"`
				PluginID               string `xml:"pluginID,attr"`
				Protocol               string `xml:"protocol,attr"`
				SvcName                string `xml:"svc_name,attr"`
				PluginModificationDate string `xml:"plugin_modification_date"`
				PluginPublicationDate  string `xml:"plugin_publication_date"`
				PluginType             string `xml:"plugin_type"`
				Solution               string `xml:"solution"`
				Description            string `xml:"description"`
				Synopsis               string `xml:"synopsis"`
				RiskFactor             string `xml:"risk_factor"`
				ScriptVersion          string `xml:"script_version"`
				PluginOutput           string `xml:"plugin_output"`
				SeeAlso                string `xml:"see_also"`
			} `xml:"ReportItem"`
		} `xml:"ReportHost"`
	} `xml:"Report"`
}
