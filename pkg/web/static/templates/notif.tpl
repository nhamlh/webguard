<!--General notification page-->
<html>
  {{ template "header" . }}
  <body>
    <div class="container">
      {{ template "display_errors" . }}
      {{ template "display_msg" . }}

      <section class="section">
      <div class="columns is-centered">
      <div class="column is-one-third">
        <a class="button is-cancel is-fullwidth" href="/">
          Go back
        </a>
      </div> <!--column-->
      </div> <!--columns-->
      </section>
    </div> <!-- container -->
  </body>
</html>
