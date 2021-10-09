
{{ define "display_msg" }}
  {{ if $.msg }}
    <br>
    <section class="section">
      <div class="columns is-centered">
        <div class="notification column is-one-third is-primary">
          {{ range $.msg }}
            <div>{{ . }}</div>
          {{ end }}
        </div>
      </div>
    </section>
  {{ end }}
{{ end }}
