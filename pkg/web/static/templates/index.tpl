<html>
  {{ template "header" }}
  <body>
    <div>My devices</div>
    <br>

    {{ range $.devices }}
      <div>{{ .name }}: {{ .key }}</div>
    {{ end }}

  </body>
</html>
