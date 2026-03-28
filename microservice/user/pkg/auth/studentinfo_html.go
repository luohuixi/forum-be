package auth

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type ccnuHTMLInput struct {
	ID    string
	Name  string
	Type  string
	Value string
}

type ccnuHTMLForm struct {
	Action string
	Method string
	Inputs []ccnuHTMLInput
}

type ccnuHTMLLink struct {
	Href string
	Text string
}

type ccnuHTMLImage struct {
	Src string
}

type ccnuPageSnapshot struct {
	CurrentURL string
	HTML       string
	Title      string
	BodyText   string
	Form       *ccnuHTMLForm
	Links      []ccnuHTMLLink
	Images     []ccnuHTMLImage
}

func parseCCNUPage(currentURL, htmlContent string) (*ccnuPageSnapshot, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	page := &ccnuPageSnapshot{
		CurrentURL: currentURL,
		HTML:       htmlContent,
		Title:      normalizeHTMLText(extractFirstElementText(doc, "title")),
		BodyText:   normalizeHTMLText(extractFirstElementText(doc, "body")),
	}

	walkHTML(doc, func(node *html.Node) {
		if node.Type != html.ElementNode {
			return
		}
		switch node.Data {
		case "form":
			if page.Form == nil {
				page.Form = parseHTMLForm(node)
			}
		case "a":
			page.Links = append(page.Links, ccnuHTMLLink{
				Href: strings.TrimSpace(getHTMLAttr(node, "href")),
				Text: normalizeHTMLText(extractNodeText(node)),
			})
		case "img":
			page.Images = append(page.Images, ccnuHTMLImage{Src: strings.TrimSpace(getHTMLAttr(node, "src"))})
		}
	})

	if page.BodyText == "" {
		page.BodyText = normalizeHTMLText(extractNodeText(doc))
	}
	return page, nil
}

func parseHTMLForm(node *html.Node) *ccnuHTMLForm {
	form := &ccnuHTMLForm{
		Action: strings.TrimSpace(getHTMLAttr(node, "action")),
		Method: strings.ToUpper(strings.TrimSpace(getHTMLAttr(node, "method"))),
	}
	if form.Method == "" {
		form.Method = "GET"
	}

	walkHTML(node, func(child *html.Node) {
		if child.Type != html.ElementNode || child.Data != "input" {
			return
		}
		form.Inputs = append(form.Inputs, ccnuHTMLInput{
			ID:    strings.TrimSpace(getHTMLAttr(child, "id")),
			Name:  strings.TrimSpace(getHTMLAttr(child, "name")),
			Type:  strings.ToLower(strings.TrimSpace(getHTMLAttr(child, "type"))),
			Value: getHTMLAttr(child, "value"),
		})
	})
	return form
}

func extractFormValues(form *ccnuHTMLForm) url.Values {
	values := url.Values{}
	if form == nil {
		return values
	}

	for _, input := range form.Inputs {
		if input.Name == "" {
			continue
		}
		switch input.Type {
		case "submit", "button", "image", "file":
			continue
		}
		values.Set(input.Name, input.Value)
	}
	return values
}

func findCaptchaURL(page *ccnuPageSnapshot) string {
	if page == nil {
		return ""
	}
	for _, image := range page.Images {
		src := strings.ToLower(image.Src)
		if strings.Contains(src, "captcha") {
			return resolveCCNUURL(page.CurrentURL, image.Src)
		}
	}
	return ""
}

func detectSecondAuthMethods(page *ccnuPageSnapshot) ([]string, string, map[string]string) {
	if page == nil {
		return nil, "", nil
	}

	methods := make([]string, 0, 2)
	seen := map[string]struct{}{}
	addMethod := func(method string) {
		if method == "" {
			return
		}
		if _, exists := seen[method]; exists {
			return
		}
		seen[method] = struct{}{}
		methods = append(methods, method)
	}

	bodyText := page.BodyText
	if strings.Contains(bodyText, "手机短信") || strings.Contains(bodyText, "手机号：") || strings.Contains(bodyText, "手机号码：") {
		addMethod(studentSecondAuthMethodSMS)
	}
	if strings.Contains(bodyText, "电子邮件") || strings.Contains(bodyText, "电子邮箱：") {
		addMethod(studentSecondAuthMethodEmail)
	}

	switchURLs := make(map[string]string, 2)
	for _, link := range page.Links {
		lowerHref := strings.ToLower(link.Href)
		text := link.Text
		if strings.Contains(text, "电子邮件") || strings.Contains(text, "电子邮箱") || strings.Contains(lowerHref, "secondloginm") {
			addMethod(studentSecondAuthMethodEmail)
			switchURLs[studentSecondAuthMethodEmail] = resolveCCNUURL(page.CurrentURL, link.Href)
		}
		if strings.Contains(text, "手机短信") || (strings.Contains(lowerHref, "secondlogin.jsf") && !strings.Contains(lowerHref, "secondloginm")) {
			addMethod(studentSecondAuthMethodSMS)
			switchURLs[studentSecondAuthMethodSMS] = resolveCCNUURL(page.CurrentURL, link.Href)
		}
	}

	currentMethod := ""
	switch {
	case strings.Contains(bodyText, "电子邮箱："):
		currentMethod = studentSecondAuthMethodEmail
	case strings.Contains(bodyText, "手机号：") || strings.Contains(bodyText, "手机号码："):
		currentMethod = studentSecondAuthMethodSMS
	case strings.Contains(strings.ToLower(page.CurrentURL), "secondloginm"):
		currentMethod = studentSecondAuthMethodEmail
	case strings.Contains(strings.ToLower(page.CurrentURL), "secondlogin"):
		currentMethod = studentSecondAuthMethodSMS
	}

	methods = orderSecondAuthMethods(methods)
	if currentMethod == "" && len(methods) == 1 {
		currentMethod = methods[0]
	}
	return methods, currentMethod, switchURLs
}

func findSecondAuthCodeField(form *ccnuHTMLForm) string {
	if form == nil {
		return ""
	}
	fallback := ""
	for _, input := range form.Inputs {
		if input.Type != "text" {
			continue
		}
		lowerName := strings.ToLower(input.Name)
		lowerID := strings.ToLower(input.ID)
		if strings.Contains(lowerName, "hidcount") || strings.Contains(lowerID, "hidcount") {
			continue
		}
		if input.Name == "" {
			continue
		}
		if strings.HasPrefix(input.Name, "form1:") {
			return input.Name
		}
		if fallback == "" {
			fallback = input.Name
		}
	}
	return fallback
}

func findSecondAuthSendButton(form *ccnuHTMLForm) (string, string, bool) {
	return findSubmitInputByValue(form, "发送验证码")
}

func findSecondAuthConfirmButton(form *ccnuHTMLForm) (string, string, bool) {
	return findSubmitInputByValue(form, "确认")
}

func findSubmitInputByValue(form *ccnuHTMLForm, keyword string) (string, string, bool) {
	if form == nil {
		return "", "", false
	}
	for _, input := range form.Inputs {
		if input.Type != "submit" {
			continue
		}
		if keyword == "" || strings.Contains(strings.TrimSpace(input.Value), keyword) {
			return input.Name, input.Value, true
		}
	}
	return "", "", false
}

func isCCNULoggedInPage(page *ccnuPageSnapshot) bool {
	if page == nil {
		return false
	}
	text := strings.ToLower(page.CurrentURL + "\n" + page.Title + "\n" + page.BodyText)
	return strings.Contains(text, "one.ccnu.edu.cn/index.html") ||
		strings.Contains(text, "one.ccnu.edu.cn/default/index.html") ||
		strings.Contains(text, "one.ccnu.edu.cn/ccnunew/index.html") ||
		strings.Contains(page.BodyText, "网上办事服务大厅") ||
		strings.Contains(page.BodyText, "一网通办")
}

func isCCNUSecondAuthPage(page *ccnuPageSnapshot) bool {
	if page == nil {
		return false
	}
	text := page.CurrentURL + "\n" + page.Title + "\n" + page.BodyText
	return strings.Contains(strings.ToLower(text), "secondlogin") || strings.Contains(text, "二次认证")
}

func isCCNUCredentialInvalidPage(page *ccnuPageSnapshot) bool {
	if page == nil {
		return false
	}
	text := page.BodyText
	return strings.Contains(text, "您输入的用户名或密码有误") ||
		strings.Contains(text, "用户名或密码有误") ||
		strings.Contains(text, "用户名或密码错误")
}

func ccnuCaptchaMessage(page *ccnuPageSnapshot) string {
	if page == nil {
		return "请输入验证码后继续登录。"
	}
	text := page.BodyText
	switch {
	case strings.Contains(text, "验证码错误"):
		return "验证码错误，请重新输入。"
	case strings.Contains(text, "验证码不正确"):
		return "验证码不正确，请重新输入。"
	default:
		return "请输入验证码后继续登录。"
	}
}

func isCCNUCaptchaChallengePage(page *ccnuPageSnapshot) bool {
	if page == nil {
		return false
	}
	text := page.Title + "\n" + page.BodyText + "\n" + page.CurrentURL
	lowerURL := strings.ToLower(page.CurrentURL)
	return findCaptchaURL(page) != "" ||
		(strings.Contains(text, "验证码") && (strings.Contains(text, "错误") || strings.Contains(text, "不正确"))) ||
		(strings.Contains(lowerURL, "cas/login") && strings.Contains(page.Title, "统一身份认证"))
}

func ccnuUnexpectedLoginMessage(page *ccnuPageSnapshot) string {
	if page == nil {
		return "登录失败，学校认证页面返回了未识别的状态，请稍后重试。"
	}
	text := page.BodyText
	switch {
	case strings.Contains(text, "系统维护"):
		return "登录失败，学校认证系统正在维护，请稍后重试。"
	case strings.Contains(text, "系统异常"), strings.Contains(text, "服务异常"):
		return "登录失败，学校认证系统返回异常页面，请稍后重试。"
	default:
		return "登录失败，原因不是验证码错误，请检查账号密码或学校站点状态。"
	}
}

func ccnuSecondAuthSendMessage(page *ccnuPageSnapshot, method string) string {
	if page == nil {
		return "验证码已发送，请输入验证码。"
	}
	text := page.BodyText
	switch {
	case strings.Contains(text, "半小时内限制发送3次验证码，请稍后重试"):
		return "半小时内限制发送3次验证码，请稍后重试！"
	case strings.Contains(text, "验证码发送成功"):
		return ccnuMethodLabel(method) + "验证码已发送，请输入验证码。"
	default:
		return ccnuMethodLabel(method) + "验证码已发送，请输入验证码。"
	}
}

func ccnuSecondAuthVerifyMessage(page *ccnuPageSnapshot) string {
	if page == nil {
		return "二次认证未通过，请重新输入验证码。"
	}
	text := page.BodyText
	switch {
	case strings.Contains(text, "验证码错误"):
		return "验证码错误，请重新输入！"
	default:
		return "二次认证未通过，请重新输入验证码。"
	}
}

func resolveCCNUURL(baseURL, rawRef string) string {
	trimmedRef := strings.TrimSpace(rawRef)
	if trimmedRef == "" {
		return baseURL
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		return trimmedRef
	}
	ref, err := url.Parse(trimmedRef)
	if err != nil {
		return trimmedRef
	}
	return base.ResolveReference(ref).String()
}

func sanitizeCaptchaText(raw string) string {
	var builder strings.Builder
	for _, r := range strings.TrimSpace(raw) {
		if (r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			builder.WriteRune(r)
		}
	}
	cleaned := builder.String()
	if len(cleaned) > 4 {
		cleaned = cleaned[len(cleaned)-4:]
	}
	return cleaned
}

func orderSecondAuthMethods(methods []string) []string {
	seen := map[string]struct{}{}
	ordered := make([]string, 0, len(methods))
	appendMethod := func(method string) {
		if method == "" {
			return
		}
		if _, exists := seen[method]; exists {
			return
		}
		seen[method] = struct{}{}
		ordered = append(ordered, method)
	}
	if containsString(methods, studentSecondAuthMethodSMS) {
		appendMethod(studentSecondAuthMethodSMS)
	}
	if containsString(methods, studentSecondAuthMethodEmail) {
		appendMethod(studentSecondAuthMethodEmail)
	}
	for _, method := range methods {
		appendMethod(method)
	}
	return ordered
}

func normalizeStudentSecondAuthMethod(method string) string {
	switch strings.ToLower(strings.TrimSpace(method)) {
	case studentSecondAuthMethodSMS:
		return studentSecondAuthMethodSMS
	case studentSecondAuthMethodEmail:
		return studentSecondAuthMethodEmail
	default:
		return ""
	}
}

func ccnuMethodLabel(method string) string {
	if method == studentSecondAuthMethodEmail {
		return "邮箱"
	}
	return "短信"
}

func walkHTML(node *html.Node, visit func(*html.Node)) {
	if node == nil {
		return
	}
	visit(node)
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		walkHTML(child, visit)
	}
}

func extractFirstElementText(node *html.Node, tag string) string {
	var text string
	walkHTML(node, func(child *html.Node) {
		if text != "" {
			return
		}
		if child.Type == html.ElementNode && child.Data == tag {
			text = extractNodeText(child)
		}
	})
	return text
}

func extractNodeText(node *html.Node) string {
	if node == nil {
		return ""
	}
	var builder strings.Builder
	var walk func(*html.Node)
	walk = func(current *html.Node) {
		if current == nil {
			return
		}
		if current.Type == html.TextNode {
			builder.WriteString(current.Data)
			builder.WriteByte(' ')
		}
		for child := current.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(node)
	return builder.String()
}

func normalizeHTMLText(text string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(text)), " ")
}

func getHTMLAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}
