package v1

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequestContactInfoTemplateButtons(t *testing.T) {
	cases := []struct {
		name    string
		buttons TemplateButtons
		json    string
	}{
		{
			name:    "value",
			buttons: TemplateButtons{Items: []Button{RequestContactInfoButton{}}},
			json:    `{"items":[{"type":"request_contact_info"}]}`,
		},
		{
			name:    "pointer",
			buttons: TemplateButtons{Items: []Button{&RequestContactInfoButton{}}},
			json:    `{"items":[{"type":"request_contact_info"}]}`,
		},
		{
			name: "existing button types",
			buttons: TemplateButtons{Items: []Button{
				&PlainButton{Label: "Reply"},
				&RequestContactInfoButton{},
				&PhoneButton{Label: "Call", Phone: "123"},
				&URLButton{Label: "Open", URL: "https://example.com"},
			}},
			json: `{"items":[{"type":"plain","label":"Reply"},{"type":"request_contact_info"},{"type":"phone","label":"Call","phone":"123"},{"type":"url","label":"Open","url":"https://example.com"}]}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.buttons)
			require.NoError(t, err)
			require.JSONEq(t, tc.json, string(data))

			var decoded TemplateButtons
			require.NoError(t, json.Unmarshal(data, &decoded))
			require.Len(t, decoded.Items, len(tc.buttons.Items))
			for i, button := range decoded.Items {
				require.Equal(t, tc.buttons.Items[i].ButtonType(), button.ButtonType())
				if button.ButtonType() == ButtonTypeRequestContactInfo {
					require.IsType(t, &RequestContactInfoButton{}, button)
				}
			}
			data, err = json.Marshal(decoded)
			require.NoError(t, err)
			require.JSONEq(t, tc.json, string(data))
		})
	}
}

func TestRequestContactInfoTemplateContracts(t *testing.T) {
	input := []byte(`{"buttons":{"items":[{"type":"request_contact_info"}]},"example":{"buttons":[[]]}}`)
	for _, target := range []any{
		&TemplateCreateWebhookData{}, &TemplateUpdateWebhookData{},
		&ActivateTemplateRequest{}, &UpdateTemplateRequest{}, &Template{},
	} {
		require.NoError(t, json.Unmarshal(input, target))
		data, err := json.Marshal(target)
		require.NoError(t, err)
		var fields map[string]json.RawMessage
		require.NoError(t, json.Unmarshal(data, &fields))
		require.JSONEq(t, `{"items":[{"type":"request_contact_info"}]}`, string(fields["buttons"]))
		require.JSONEq(t, `{"buttons":[[]]}`, string(fields["example"]))
	}

	var templates []Template
	require.NoError(t, json.Unmarshal(append(append([]byte("["), input...), ']'), &templates))
	require.Len(t, templates, 1)
	require.IsType(t, &RequestContactInfoButton{}, templates[0].Buttons.Items[0])
	require.Equal(t, [][]string{{}}, templates[0].Example.Buttons)
}

func TestRequestContactInfoSuggestion(t *testing.T) {
	suggestion := Suggestion{Type: SuggestionTypeRequestContactInfo}
	data, err := json.Marshal(suggestion)
	require.NoError(t, err)
	require.JSONEq(t, `{"type":"request_contact_info"}`, string(data))
	var decoded Suggestion
	require.NoError(t, json.Unmarshal(data, &decoded))
	require.Equal(t, suggestion, decoded)
}

func TestRequestContactInfoChannelSettings(t *testing.T) {
	settings := ChannelSettingsSuggestions{Text: ChannelFeatureBoth, RequestContactInfo: ChannelFeatureReceive}
	data, err := json.Marshal(settings)
	require.NoError(t, err)
	require.JSONEq(t, `{"text":"both","request_contact_info":"receive"}`, string(data))
	var decoded ChannelSettingsSuggestions
	require.NoError(t, json.Unmarshal(data, &decoded))
	require.Equal(t, settings, decoded)

	data, err = json.Marshal(ChannelSettingsSuggestions{Text: ChannelFeatureBoth})
	require.NoError(t, err)
	require.JSONEq(t, `{"text":"both"}`, string(data))
}
