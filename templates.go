package main

import (
	"html/template"
	"io"
	"sort"
	"time"

	"github.com/zulily/stevedore/repo"
)

const (
	htmlTmpl = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>stevedore</title>

    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.4/css/bootstrap.min.css">
		<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/highlight.js/8.5/styles/default.min.css">
  </head>
  <body>

	<nav class="navbar navbar-default">
		<div class="container-fluid">
			<div class="navbar-header">
				<a class="navbar-brand" href="#">
					<p>ᕕ( ⁰ ▽ ⁰ )ᕗ </p>
				</a>

				<form id="add-repo" class="navbar-form navbar-left">
					<div class="form-group">
						<input name="repo" type="text" class="form-control" style="width: 500px; important!" placeholder="https://github.com/username/project.git">
					</div>
					<button id="submit" type="submit" class="btn btn-default">Add new repo</button>
				</form>
			</div>
		</div>
	</nav>

    <div class="container">
			<div id="result"></div>
      <div class="jumbotron">
        <h1>stevedore</h1>
				<p>ᕕ( ⁰ ▽ ⁰ )ᕗ  buildin' your containers</p>
      </div>

        {{ range . }}
				<div class="panel panel-primary">
					<div class="panel-heading"><h3>{{ .URL }}</h3></div>
					<div class="panel-body">
						<ul class="list-group">

						  {{ if .SHA }}
							<li class="list-group-item">
							  SHA:<span class="label label-primary pull-right">{{ .SHA }}</span>
							</li>
							{{ end }}

							<li class="list-group-item">
								Build status:
								{{ if isFailing . }}
								<span class="label label-danger pull-right">Failing</span>
								{{ else if isPassing . }}
								<span class="label label-success pull-right">Passing</span>
								{{ else if isInProgress . }}
								<span class="label label-warning pull-right">In progress</span>
								{{ end }}
							</li>

							{{ if .LastPublishDate }}
							<li class="list-group-item">
							  Last published:<span class="label label-primary pull-right">{{ fmtTime .LastPublishDate }}</span>
							</li>
							{{ end }}

						</ul>

            {{ if .Images }}
						<h4>Images:</h4>
						<ul class="list-group">
						  {{ range $img := .Images }}
							  <li class="list-group-item">{{ $img }}</li>
							{{ end }}
						</ul>
						{{ end }}

				  	{{ if .Log }}
 						<h4>Last build message:</h4>
						<pre><code class="bash">
{{ .Log }}
						</code></pre>
						{{ end }}
					</div>
				</div>
				{{ end }}

     </container>
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.2/jquery.min.js"></script>
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.4/js/bootstrap.min.js"></script>
		<script src="//cdnjs.cloudflare.com/ajax/libs/highlight.js/8.5/highlight.min.js"></script>
		<script>hljs.initHighlightingOnLoad();</script>
		<script type="text/javascript">
			$(document).ready(function () {
				$("#add-repo").submit(function(e){
					e.preventDefault();
					var repoRequest = {};
					repoRequest["repo"] = $( "#add-repo :input[name=repo]" ).val();

					$.ajax({
						type: "POST",
						url: "/repos",
						data: JSON.stringify(repoRequest),
						success: function(msg){
							$("#result").html('<div class="alert alert-success"><button type="button" class="close">×</button>Added a new repo!</div>');
							window.setTimeout(function() {
								$(".alert").fadeTo(500, 0).slideUp(500, function(){
									$(this).remove();
								});
							}, 1000);
						},

						error: function(xhr, status, err){
							$("#result").html('<div class="alert alert-danger"><button type="button" class="close">×</button>' + err + '</div>');
							window.setTimeout(function() {
								$(".alert").fadeTo(500, 0).slideUp(500, function(){
									$(this).remove();
								});
							}, 1000);
						}

					}); // ajax
				}); 	// form submit
			});
		</script>
	</body>
</html>
`
)

type Repos []*repo.Repo

// FormattedLastPublishDate returns a human-readable, formatted representation of a
// unix epoch-seconds timestamp.
func formatUnixTime(unixTime int64) string {
	if unixTime == 0 {
		return ""
	}
	return time.Unix(unixTime, 0).Format(time.RFC1123)
}

func isFailing(r *repo.Repo) bool {
	return r.Status == repo.StatusFailing
}

func isInProgress(r *repo.Repo) bool {
	return r.Status == repo.StatusInProgress
}

func isPassing(r *repo.Repo) bool {
	return r.Status == repo.StatusPassing
}

func (slice Repos) Len() int {
	return len(slice)
}

// Less sorts "in progress" repos to the top, followed by the rest of the Repos
// by publish date (descending), then by name (lexicographically)
func (slice Repos) Less(i, j int) bool {
	if slice[i].Status == repo.StatusInProgress && slice[j].Status != repo.StatusInProgress {
		return true
	}

	if slice[j].Status == repo.StatusInProgress && slice[i].Status != repo.StatusInProgress {
		return false
	}

	if slice[i].LastPublishDate != slice[j].LastPublishDate {
		return slice[i].LastPublishDate > slice[j].LastPublishDate
	}

	return slice[i].URL < slice[j].URL
}

func (slice Repos) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// Render renders a Services instance to the supplied io.Writer as HTML
func RenderServicesHTML(repos Repos, out io.Writer) error {
	sort.Sort(repos)
	funcMap := template.FuncMap{
		"fmtTime":      formatUnixTime,
		"isPassing":    isPassing,
		"isFailing":    isFailing,
		"isInProgress": isInProgress,
	}
	t := template.Must(template.New("htmloutput").Funcs(funcMap).Parse(htmlTmpl))
	return t.Execute(out, repos)
}
