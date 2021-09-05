{{ define "header" }}
<header>
  <link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.12.1/css/all.css" crossorigin="anonymous">
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@0.9.3/css/bulma.min.css">
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma-divider@0.2.0/dist/css/bulma-divider.min.css">

  <title>Simple Wireguard dashboard</title>

  <nav class="navbar is-primary is-fluid" role="navigation" aria-label="main navigation">
    <div class="navbar-brand">
      <a class="navbar-item" href="#">
        <img src="#" alt="Webguard">
      </a>

      <a role="button" class="navbar-burger" aria-label="menu" aria-expanded="false" data-target="navbarBasicExample">
        <span aria-hidden="true"></span>
        <span aria-hidden="true"></span>
        <span aria-hidden="true"></span>
      </a>
    </div>

    <div id="navbar" class="navbar-menu">
      <div class="navbar-start">
        <a class="navbar-item" href="/">
          <span class="icon-text">
            <span class="icon">
              <i class="fas fa-home"></i>
            </span>
            <span>Home</span>
          </span>
        </a>
      </div>

      <div class="navbar-end">
        <div class="navbar-item">
          <div class="buttons">
            <a href="/login" class="button is-light">
              <span class="icon-text">
                <span class="icon">
                  <i class="fas fa-sign-in-alt"></i>
                </span>
                <span>Login</span>
              </span>
            </a>
          </div>
        </div>
      </div>
    </div>
  </nav>
</header>
{{ end  }}
