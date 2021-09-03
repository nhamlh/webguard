<html>
  {{ template "header" }}
  <body class="has-background-info-light">
    <div class="container">

      {{ template "display_errors" . }}

      <section class="section is-small">
        <form class="box" action="/login" method="POST">
          <div class="field">
            <label class="label">Email</label>
            <div class="control">
              <input class="input" type="text" name="email" placeholder="e.g. bob@example.com">
            </div>
          </div>

          <div class="field">
            <label class="label">Password</label>
            <div class="control">
              <input class="input" type="password" name="password" placeholder="********">
            </div>
          </div>

          <div class="field is-grouped">
            <div class="control">
              <button class="button is-primary">Submit</button>
            </div>
            <div class="control">
              <button class="button is-link is-light">Cancel</button>
            </div>
          </div>
        </form>

        <hr>

        <a class="button is-primary" href="/login/oauth">
          <span class="icon">
            <i class="fas fa-github"></i>
          </span>
          <span>Login with GitHub</span>
        </a>
      </section>

    </div> <!-- container -->
  </body>
</html>
