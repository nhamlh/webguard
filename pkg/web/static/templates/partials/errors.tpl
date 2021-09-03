
{{ define "display_errors" }}
  {{ if $.errors }}
    <br>
    <section class="notification section is-danger">
      {{ range $.errors }}
        <div>{{ . }}</div>
      {{ end }}
    </section>
  {{ end }}
{{ end }}
