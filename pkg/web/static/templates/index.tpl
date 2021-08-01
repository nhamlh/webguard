<html>
  {{ template "header" }}
  <body>
    <div>My devices</div>
    <br>

    {{ range $.devices }}
      <div>{{ .Name }}: {{ .AllowedIps }}</div>
    {{ end }}

  </body>
</html>
