<html>
  {{ template "header" . }}
  <body>
    <div class="container">
      {{ template "display_errors" . }}

      <section class="section">
      <div class="columns is-centered">
      <div class="column is-one-third has-background-light">

        <form action="/new_device" method="POST">
          <div class="field">
            <label class="label">Name</label>
            <div class="control has-icons-left">
              <input class="input" type="text" name="name" placeholder="e.g. My Mac Device">
              <span class="icon is-left">
                <i class="fas fa-laptop-code"></i>
              </span>
            </div>
          </div>

          <div class="field is-grouped">
            <div class="control">
              <button class="button is-primary">
                <span class="icon">
                  <i class="fas fa-plus"></i>
                </span>
                <span>Create</span>
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
