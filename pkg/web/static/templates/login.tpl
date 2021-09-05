<html>
  {{ template "header" . }}
  <body>
    <div class="container">

      {{ template "display_errors" . }}

      <section class="section">
      <div class="columns is-centered">
      <div class="column is-one-third has-background-light">

        <div class="title is-5 has-text-centered">
              Single Sign-On
          </span>
        </div>

        <a class="button is-primary is-fullwidth" href="/login/oauth">
          <span class="icon">
            <i class="fas fa-sign-in-alt"></i>
          </span>
          <span>SSO Login</span>
        </a>
      </div> <!--column-->
      </div> <!--columns-->
      <div class="columns is-centered">
      <div class="column is-one-third has-background-light">
        <div class="is-divider"></div>

        <div class="title is-5 has-text-centered">
              Static Login
          </span>
        </div>

        <form action="/login" method="POST">
          <div class="field">
            <label class="label">Email</label>
            <div class="control has-icons-left">
              <input class="input" type="text" name="email" placeholder="e.g. bob@example.com">
              <span class="icon is-left">
                <i class="fas fa-envelope"></i>
              </span>
            </div>
          </div>

          <div class="field">
            <label class="label">Password</label>
            <div class="control has-icons-left">
              <input class="input" type="password" name="password" placeholder="********">
              <span class="icon is-left">
                <i class="fas fa-lock"></i>
              </span>
            </div>
          </div>

          <div class="control">
            <button class="button is-primary is-fullwidth">
              <span class="icon">
                <i class="fas fa-sign-in-alt"></i>
              </span>
              <span>Login</span>
            </button>
          </div>
        </form>
      </div> <!--column-->
      </div> <!--columns-->
      </section>

    </div> <!-- container -->
  </body>
</html>
