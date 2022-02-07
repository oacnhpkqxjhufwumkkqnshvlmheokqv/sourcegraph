// Package suspiciousnames contains a blocklist of suspicious names for users and organizations.
package suspiciousnames

import (
	"strings"

	"github.com/sourcegraph/sourcegraph/cmd/frontend/envvar"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

// CheckNameAllowedForUserOrOrganization returns a non-nil error if the desired username or
// organization name is suspicious and the currently running site is Sourcegraph.com.
//
// It is intended to prevent users/organizations on Sourcegraph.com from changing their name to a
// name that is likely to be confused with an admin or official account (such as "admin", "root",
// etc.).
//
// 🚨 SECURITY: This is not foolproof; users may choose a name like `secur1ty` that might be
// confused with a name like "security", or they might find another synonym that we didn't think of.
func CheckNameAllowedForUserOrOrganization(desiredName string) error {
	if envvar.SourcegraphDotComMode() && isSuspicious(desiredName) {
		return errors.Errorf("rejected suspicious name %q", desiredName)
	}
	return nil
}

func isSuspicious(desiredName string) bool {
	_, ok := suspiciousNames[strings.ToLower(desiredName)]
	return ok
}

var suspiciousNames = map[string]struct{}{
	"about":          {},
	"access":         {},
	"account":        {},
	"accounts":       {},
	"add":            {},
	"address":        {},
	"adm":            {},
	"admin":          {},
	"administration": {},
	"adult":          {},
	"advertising":    {},
	"affiliate":      {},
	"affiliates":     {},
	"ajax":           {},
	"analytics":      {},
	"android":        {},
	"anon":           {},
	"anonymous":      {},
	"api":            {},
	"app":            {},
	"apps":           {},
	"archive":        {},
	"atom":           {},
	"auth":           {},
	"authentication": {},
	"avatar":         {},
	"backup":         {},
	"banner":         {},
	"banners":        {},
	"billing":        {},
	"bin":            {},
	"blog":           {},
	"blogs":          {},
	"board":          {},
	"bot":            {},
	"bots":           {},
	"business":       {},
	"cache":          {},
	"calendar":       {},
	"campaign":       {},
	"careers":        {},
	"cgi":            {},
	"chat":           {},
	"client":         {},
	"code":           {},
	"commercial":     {},
	"compare":        {},
	"config":         {},
	"connect":        {},
	"contact":        {},
	"contest":        {},
	"cpanel":         {},
	"create":         {},
	"css":            {},
	"dashboard":      {},
	"data":           {},
	"database":       {},
	"delete":         {},
	"demo":           {},
	"design":         {},
	"designer":       {},
	"dev":            {},
	"devel":          {},
	"dir":            {},
	"directory":      {},
	"doc":            {},
	"docs":           {},
	"domain":         {},
	"download":       {},
	"downloads":      {},
	"ecommerce":      {},
	"edit":           {},
	"editor":         {},
	"email":          {},
	"example":        {},
	"explore":        {},
	"faq":            {},
	"favorite":       {},
	"feed":           {},
	"feedback":       {},
	"file":           {},
	"files":          {},
	"flog":           {},
	"follow":         {},
	"forum":          {},
	"forums":         {},
	"free":           {},
	"ftp":            {},
	"games":          {},
	"group":          {},
	"groups":         {},
	"guest":          {},
	"help":           {},
	"home":           {},
	"homepage":       {},
	"host":           {},
	"hosting":        {},
	"hostname":       {},
	"html":           {},
	"http":           {},
	"httpd":          {},
	"https":          {},
	"image":          {},
	"images":         {},
	"imap":           {},
	"img":            {},
	"index":          {},
	"info":           {},
	"information":    {},
	"intranet":       {},
	"invite":         {},
	"ipad":           {},
	"iphone":         {},
	"irc":            {},
	"issues":         {},
	"java":           {},
	"javascript":     {},
	"job":            {},
	"jobs":           {},
	"js":             {},
	"knowledgebase":  {},
	"list":           {},
	"lists":          {},
	"log":            {},
	"login":          {},
	"logout":         {},
	"logs":           {},
	"mail":           {},
	"mail1":          {},
	"mail2":          {},
	"mail3":          {},
	"mail4":          {},
	"mail5":          {},
	"mailer":         {},
	"mailing":        {},
	"manager":        {},
	"marketing":      {},
	"marketplace":    {},
	"master":         {},
	"me":             {},
	"media":          {},
	"message":        {},
	"messenger":      {},
	"mine":           {},
	"mobile":         {},
	"movie":          {},
	"movies":         {},
	"msg":            {},
	"msn":            {},
	"music":          {},
	"mx":             {},
	"my":             {},
	"mysql":          {},
	"name":           {},
	"named":          {},
	"net":            {},
	"network":        {},
	"new":            {},
	"news":           {},
	"newsletter":     {},
	"nick":           {},
	"nickname":       {},
	"notes":          {},
	"ns":             {},
	"ns1":            {},
	"ns2":            {},
	"ns3":            {},
	"ns4":            {},
	"old":            {},
	"online":         {},
	"operator":       {},
	"order":          {},
	"orders":         {},
	"page":           {},
	"pager":          {},
	"pages":          {},
	"panel":          {},
	"password":       {},
	"perl":           {},
	"photo":          {},
	"photos":         {},
	"php":            {},
	"pic":            {},
	"pics":           {},
	"plugin":         {},
	"plugins":        {},
	"pop":            {},
	"pop3":           {},
	"post":           {},
	"postfix":        {},
	"postmaster":     {},
	"posts":          {},
	"profile":        {},
	"project":        {},
	"projects":       {},
	"promo":          {},
	"pub":            {},
	"public":         {},
	"pulls":          {},
	"python":         {},
	"random":         {},
	"register":       {},
	"registration":   {},
	"root":           {},
	"rss":            {},
	"ruby":           {},
	"sale":           {},
	"sales":          {},
	"sample":         {},
	"samples":        {},
	"script":         {},
	"scripts":        {},
	"search":         {},
	"secure":         {},
	"security":       {},
	"send":           {},
	"service":        {},
	"setting":        {},
	"settings":       {},
	"setup":          {},
	"shop":           {},
	"signin":         {},
	"signup":         {},
	"site":           {},
	"sitemap":        {},
	"sites":          {},
	"smtp":           {},
	"sql":            {},
	"ssh":            {},
	"staff":          {},
	"stafftools":     {},
	"stage":          {},
	"staging":        {},
	"start":          {},
	"stat":           {},
	"static":         {},
	"stats":          {},
	"status":         {},
	"store":          {},
	"stores":         {},
	"subdomain":      {},
	"subscribe":      {},
	"suporte":        {},
	"support":        {},
	"system":         {},
	"talk":           {},
	"task":           {},
	"tasks":          {},
	"tech":           {},
	"telnet":         {},
	"test":           {},
	"test1":          {},
	"test2":          {},
	"test3":          {},
	"teste":          {},
	"tests":          {},
	"theme":          {},
	"themes":         {},
	"tmp":            {},
	"todo":           {},
	"tools":          {},
	"update":         {},
	"upload":         {},
	"url":            {},
	"usage":          {},
	"user":           {},
	"username":       {},
	"video":          {},
	"videos":         {},
	"visitor":        {},
	"web":            {},
	"webmail":        {},
	"webmaster":      {},
	"website":        {},
	"websites":       {},
	"ww":             {},
	"www":            {},
	"www1":           {},
	"www2":           {},
	"www3":           {},
	"www4":           {},
	"www5":           {},
	"www6":           {},
	"www7":           {},
	"wwws":           {},
	"wwww":           {},
	"xpg":            {},
	"xxx":            {},
	"you":            {},
	"yourdomain":     {},
	"yourname":       {},
	"yoursite":       {},
	"yourusername":   {},
}
