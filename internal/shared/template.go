package shared

import (
	"bufio"
	"bytes"
	"html/template"
	"time"
)

func CloudDiscoveryPage(result CloudDiscoveryResults) []byte {
	body, err := resolveTemplate(`
          <table width="100%" border="0" cellspacing="0" cellpadding="0">
            <thead style="background-color: #E2E6F9; color: #9FACB5; letter-spacing: 1px;">
				<tr style="background-color: #4e1fcc; text-transform: none; font-size: 18px;">
				<th style="padding: 5px;"><img src="logo.png" style="width: 200px;"/></th></tr>
				<tr style="background-color: #4e1fcc; color: white; text-transform: none; font-size: 18px;">
				<th style="padding-top: 0px; padding-bottom: 15px;">AWS Cloud Discovery Results</th></tr>
				<tr style="background-color: #4e1fcc; color: white; text-transform: none; font-size: 11px;">
				<th style="padding: 5px;">For more info <a href="http://www.twistlock.com" style="color: white; cursor: pointer">click here</a></th></tr>
			</thead>
			<tbody>
				{{ range .Results }}
						<tr><td style="height: 15px"></td></tr>
						<tr>
							<td style="background-color: #E2E6F9; text-align: center; font-size: 17px; padding: 5px 15px;">
								<strong>Type: <a style="color: #000000">{{ .Type }}</a></strong>
							</td>
						</tr>
						<tr>
							<td style="background-color: #E2E6F9; text-align: center; font-size: 17px; padding: 5px 15px;">
								<strong>Region: <a style="color: #000000">{{ .Region }}</a></strong>
							</td>
						</tr>
						<tr>
							<td>
								<table width="100%" border="1px" cellspacing="0" cellpadding="0">
									<tr>
										<td style="font-size: 16px; padding: 5px 15px; border: none">
											<strong>Discoveries:</strong>
										</td>
									</tr>
									{{ range .Assets}}
									<tr><td style="font-size: 14px; padding: 0px 15px; border: none">{{ .ID }}</td></tr>
									{{ end }}
								</table>
							</td>
						</tr>
						<tr><td style="height: 2px"></td></tr>
				{{ end }}
			</tbody>
          </table>
    `, struct {
		Time     string
		Results []CloudDiscoveryResult
	}{
		Time:     time.Now().Format(time.RFC822Z),
		Results: result.Results,
	})
	if err != nil {
		panic(err)
	}
	return body
}

// ResolveTemplate resolves the given string using templates formatting
// https://golang.org/pkg/text/template/
func resolveTemplate(tmplt string, data interface{}) ([]byte, error) {
	t, err := template.New("template").Parse(tmplt)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	if err := t.ExecuteTemplate(writer, "template", data); err != nil {
		return nil, err
	}

	writer.Flush()
	return buf.Bytes(), nil
}