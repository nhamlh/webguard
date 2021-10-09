<html>
  {{ template "header" . }}
  <body>
    <div class="container">
      {{ template "display_errors" . }}
      {{ template "display_msg" . }}

      <section class="section">
      <div class="columns is-centered">
      <div class="column is-one-third">
        <form action="/change_password" method="POST">
          <div class="field">
            <label class="label">Your current password</label>
            <div class="control has-icons-left">
              <input class="input" type="password" name="current_password" placeholder="********">
              <span class="icon is-left">
                <i class="fas fa-lock"></i>
              </span>
            </div>
          </div>

          <div class="field">
            <label class="label">Your new password</label>
            <div class="control has-icons-left">
              <input class="input" type="password" name="new_password" placeholder="********">
              <span class="icon is-left">
                <i class="fas fa-lock"></i>
              </span>
            </div>
          </div>

          <div class="field is-grouped">
            <div class="control">
              <button class="button is-primary">
                <span class="icon">
                  <i class="fas fa-check"></i>
                </span>
                <span>Change</span>
              </button>
            </div>
          </div>
        </form>
      </div> <!--column-->
      </div> <!--columns-->
      </section>
    </div> <!-- container -->
  </body>
</html>
