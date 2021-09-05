
{{ define "display_errors" }}
  {{ if $.errors }}
    <br>
    <section class="section">
      <div class="columns is-centered">
        <div class="notification column is-one-third is-warning">
          {{ range $.errors }}
            <div>{{ . }}</div>
          {{ end }}
        </div>
      </div>
    </section>
  {{ end }}
{{ end }}
