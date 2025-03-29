// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.857
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "senatus/src/db/repo"
import "strconv"

func Index(timeSlots []repo.TimeSlotModel) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>Senatus Areae Localis Nexus</title><link href=\"/static/bootstrap.min.css\" rel=\"stylesheet\"><link rel=\"stylesheet\" href=\"/static/style.css\"><script src=\"/static/htmx.js\"></script><link rel=\"icon\" href=\"/static/favicon.ico\" type=\"image/x-icon\"></head><body><!-- Header with Roman-inspired design --><header class=\"roman-header py-4 shadow\"><div class=\"container\"><div class=\"row align-items-center\"><div class=\"col-md-8 mb-3 mb-md-0 d-flex align-items-center\"><h1 class=\"mb-0\">SENATVS AREAE LOCALIS NEXVS</h1></div><div class=\"col-md-4 d-flex justify-content-md-end\"><button class=\"roman-btn btn me-2\">FORVM</button> <button class=\"roman-btn btn\">CIVES</button></div></div></div></header><!-- Main content --><div class=\"container py-5\"><div class=\"card roman-border\"><div class=\"card-header roman-red-bg text-center text-white py-3\"><h2 class=\"mb-1\">CONSVLTATIO POPVLI</h2><p class=\"mb-0 fst-italic\">The People's Consultation</p></div><div class=\"card-body p-4\"><p class=\"text-center mb-4 text-secondary\">Citizens of the Senate, cast your votes on the proposed activities.  The will of the people shall determine our course of action.</p><!-- Add new time slot form at the top --><div class=\"card roman-border mb-5\"><div class=\"card-header roman-red-bg text-white d-flex justify-content-between align-items-center\"><h3 class=\"mb-0\">NEW TIME SLOT</h3><span class=\"badge roman-gold-bg roman-red-text\">PROPOSITIO NOVA</span></div><div class=\"card-body\"><form hx-post=\"/\" hx-target=\"body\" hx-swap=\"outerHTML\"><div class=\"mb-3\"><label for=\"time\" class=\"form-label\">Select Time:</label> <input type=\"time\" id=\"time\" name=\"time\" class=\"form-control\" required></div><div class=\"d-flex justify-content-end\"><button type=\"submit\" class=\"roman-btn btn\">Create Time Slot</button></div></form></div></div><!-- Time slots container --><div id=\"timeslots-container\"><!-- 20:00 Time Slot -->")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, timeSlot := range timeSlots {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "<div class=\"card roman-border mb-4\"><div class=\"card-header roman-red-bg text-white d-flex justify-content-between align-items-center\"><h3 class=\"mb-0\">HORA ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var2 string
			templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(timeSlot.Time)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/templates/index.templ`, Line: 70, Col: 47}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "</h3><span class=\"badge roman-gold-bg roman-red-text\">5 ACTIVITIES</span></div><!-- Add new activity form at the top --><div class=\"card-body bg-light border-bottom\"><form hx-post=\"/activities\" hx-target=\"body\" hx-swap=\"outerHTML\"><input type=\"hidden\" name=\"timeSlot\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(strconv.FormatInt(timeSlot.ID, 10))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/templates/index.templ`, Line: 76, Col: 90}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "\"><div class=\"input-group\"><input type=\"text\" name=\"activity\" class=\"form-control\" placeholder=\"Propose a new activity...\" required> <button type=\"submit\" class=\"roman-btn btn\">PROPOSE</button></div></form></div><ul id=\"activities-20\" class=\"list-group list-group-flush\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, activity := range timeSlot.Activities {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "<!-- Leading activity --> <li class=\"list-group-item leading-activity d-flex align-items-center\"><div class=\"me-3 d-flex flex-column align-items-center\"><button class=\"vote-btn\" hx-post=\"/api/vote\" hx-vals=\"{&#34;timeSlotId&#34;: 1, &#34;activityId&#34;: 3, &#34;increment&#34;: true}\" hx-target=\"#activities-20\" hx-swap=\"innerHTML\">▲</button> <span class=\"vote-count my-1\">15</span> <button class=\"vote-btn\" hx-post=\"/api/vote\" hx-vals=\"{&#34;timeSlotId&#34;: 1, &#34;activityId&#34;: 3, &#34;increment&#34;: false}\" hx-target=\"#activities-20\" hx-swap=\"innerHTML\">▼</button></div><div class=\"flex-grow-1\"><span class=\"activity-name\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var4 string
				templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(activity.Name)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/templates/index.templ`, Line: 93, Col: 56}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "</span></div><div class=\"ms-auto d-flex align-items-center text-warning\"><svg xmlns=\"http://www.w3.org/2000/svg\" width=\"16\" height=\"16\" fill=\"currentColor\" class=\"bi bi-award me-1\" viewBox=\"0 0 16 16\"><path d=\"M9.669.864 8 0 6.331.864l-1.858.282-.842 1.68-1.337 1.32L2.6 6l-.306 1.854 1.337 1.32.842 1.68 1.858.282L8 12l1.669-.864 1.858-.282.842-1.68 1.337-1.32L13.4 6l.306-1.854-1.337-1.32-.842-1.68L9.669.864zm1.196 1.193.684 1.365 1.086 1.072L12.387 6l.248 1.506-1.086 1.072-.684 1.365-1.51.229L8 10.874l-1.355-.702-1.51-.229-.684-1.365-1.086-1.072L3.614 6l-.25-1.506 1.087-1.072.684-1.365 1.51-.229L8 1.126l1.356.702 1.509.229z\"></path> <path d=\"M4 11.794V16l4-1 4 1v-4.206l-2.018.306L8 13.126 6.018 12.1 4 11.794z\"></path></svg> <span class=\"fw-bold small\">LEADING</span></div></li>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, "</ul></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, "</div></div></div></div><!-- Footer --><footer class=\"roman-red-bg text-white py-3 mt-5\"><div class=\"container text-center\"><p class=\"mb-0\">SENATVS POPVLVSQVE ROMANVS • <span id=\"current-year\"></span></p></div></footer><!-- Bootstrap JS Bundle with Popper --><script src=\"https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha1/dist/js/bootstrap.bundle.min.js\"></script><!-- Set current year in footer --><script>\n    document.getElementById('current-year').textContent = new Date().getFullYear();\n  </script></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
