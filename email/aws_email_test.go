package email

import (
	"testing"
	"time"
)

func init() {
	Init("us-east-1", "no-reply@vesync.com", "AKIAJSJFZLXSOQS6O56A", "gaplcwPEa/LCqC4DJjMHi5UkujGzn6PgxtwZLGAb")
}

func TestSendEmail(t *testing.T) {
	err := sendEmail("tabzhang@etekcity.com.cn", "test", "this is a test email")
	if err != nil {
		t.Error(err)
		return
	}
}

var baseTemplate = `
<html>

<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <style class="fox_global_style">
        div.fox_html_content {
            line-height: 1.5;
        }

        /* 一些默认样式 */
        blockquote {
            margin-Top: 0px;
            margin-Bottom: 0px;
            margin-Left: 0.5em
        }

        ol,
        ul {
            margin-Top: 0px;
            margin-Bottom: 0px;
            list-style-position: inside;
        }

        p {
            margin-Top: 0px;
            margin-Bottom: 0px
        }
    </style> <!-- StartSystemHeader -->
    <!-- EndSystemheader -->
</head>

<body>
<table width="100%" border="0" cellspacing="0" cellpadding="0" class="email-body-wrapper" style="border-collapse: collapse;">
    <tbody>
    <tr>
        <td>
            <table cellspacing="0" cellpadding="0" border="0" bordercollapse="collapse" align="center" width="620"
                   id="sc3405" style="table-layout: auto; empty-cells: show; border-collapse: collapse; background-color: rgb(255, 255, 255);">
                <tbody>
                <tr>
                    <td valign="top" align="left" rowspan="1" colspan="3" width="620" height="80" id="view0">
                        <div id="sc5541" class="sc-view" style="left: 0px; width: 620px; top: 0px; height: 80px; overflow: hidden">
                            <div class="co-border-style" style="border-width: 2px; border-style: none">
                                <table width="620" height="80" cellspacing="0" cellpadding="0" border="0"
                                       bordercollapse="collapse" class="co-style-table" style="margin: 0px; border-collapse: collapse;">
                                    <tbody>
                                    <tr>
                                        <td valign="top" class="valign-able">
                                            <span style="display: block; border: 0px; border-image-source: initial; border-image-slice: initial; border-image-width: initial; border-image-outset: initial; border-image-repeat: initial; outline: none; text-decoration: none;">VeSync&nbsp;|&nbsp;云平台</span>
                                        </td>
                                    </tr>
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </td>
                </tr>
                <tr>
                    <td align="left" valign="top" width="45" height="10" id="empty3"></td>
                </tr>
                <tr>
                    <td align="left" valign="top" width="45" height="300" id="empty6"></td>
                    <td valign="top" align="left" rowspan="1" colspan="1" width="530" id="view7"
                        style="color: #000000; font-family: Arial; font-size: 12px; line-height: 18px; letter-spacing: 0px; word-wrap: break-word">
                        <div id="sc5554" class="sc-view hidden-border inline-styled-view editor-outline"
                             style="left: 45px; width: 530px; top: 90px; color: #000000; font-family: Arial; font-size: 12px; line-height: 18px; letter-spacing: 0px; word-wrap: break-word; overflow: hidden">
                            <div class="co-border-style" style="">
                                <table width="530" cellspacing="0" cellpadding="0" border="0"
                                       bordercollapse="collapse" class="co-style-table" style="color: rgb(0, 0, 0); font-family: Arial; font-size: 12px; line-height: 18px; letter-spacing: 0px; word-wrap: break-word; margin: 0px; border-collapse: collapse;">
                                    <tbody>
                                    <tr>
                                        <td valign="top" class="valign-able">
                                            <span class="remove-absolute">
                                                <span style="font-family:&#39;Akzidenz-Grotesk Std&#39;,&#39;Helvetica Neue&#39;, Helvetica, Arial, sans-serif;font-size:14px;line-height:140%;color:rgb(49, 48, 47);">
{{template "email-content" .}}
                                                    <br><br> 恭祝愉快！
                                                    <br><br> VeSync 团队
                                                </span>
                                            </span>
                                        </td>
                                    </tr>
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </td>
                    <td align="left" valign="top" width="45" height="400" id="empty8"></td>
                </tr>
                </tbody>
            </table>
        </td>
    </tr>
    </tbody>
</table>
<table align="center" height="157" style="width: 620px; height: 157px; background-color: #ffffff; color: #616161;border-collapse: collapse;"
       width="620">
    <tbody>
    <tr height="15" style="height: 15px;">
        <td colspan="6">&nbsp;</td>
    </tr>
    <tr height="44" style="height: 44 px;">
        <td align="center" colspan="6" style="vertical-align:middle;" valign="middle"><span style="font-family:&#39;Akzidenz-Grotesk Std&#39;,&#39;Helvetica Neue&#39;, Helvetica, Arial, sans-serif;font-size: 8px;line-height:125%;color:#616161;">VeSync
                        Inc., 2018. All rights reserved.<br>
                        F10, building 10, tianjie times, longhu times, oil road, yuzhong district, chongqing<br>
                        重庆市渝中区龙湖时代天街D馆龙湖时代汇10栋10层 400042 <br>
                        <br>
                        <a href="#" style="font-size:10px;color:#69b241; text-decoration:none;">系统邮件，请勿回复</a> </span></td>
    </tr>
    </tbody>
</table>
<img src="#" alt="" border="0" width="1px" height="1px" style="border: 0px; border-image-source: initial; border-image-slice: initial; border-image-width: initial; border-image-outset: initial; border-image-repeat: initial; outline: none; text-decoration: none;">
</body>

</html>
`
var content = `
{{define  "email-content"}}
<p style="margin: 1em 0px;">
    {{.ToName}},
</p>
<br>
<p style="margin: 1em 0px;">
{{.UserName}} 在我平台创建的{{.ProductName}}，在{{.Time}}提交了发布申请，请您及时审核处理。
</p>
<p style="margin: 1em 0px;">
用户名：{{.UserName}}<br>
产品名称：{{.ProductName}}<br>
产品Pid：{{.Pid}}<br>   
</p>
<p style="margin: 1em 0px;">

接下来，与用户确认<br>
1.软硬件测试是否通过，如有必要请Check测试报告<br>
2.正式版本的固件是否已经上传并确认版本是否正确<br>
3.完成所有检查项后，发布上线该产品<br>
4.通知用户已经发布完成<br>
5.通知平台维护人员该产品已经发布，线上跟踪测试<br>
<br>
</p>
{{end}}
`

func TestSendEmailWithTemplate(t *testing.T) {
	l, _ := time.LoadLocation("")
	err := SendEmailWithTemplate("tabzhang@etekcity.com.cn", "test", baseTemplate, content, map[string]interface{}{
		"ToName":      "tab",
		"ProductName": "wifi插座",
		"Pid":         "sdsdsdsdsdsdsdsd",
		"Time":        time.Now().In(l).Format(time.UnixDate),
		"UserNmae":    "Tim",
	})
	if err != nil {
		t.Error(err)
		return
	}
}
