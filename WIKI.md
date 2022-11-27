**How to do update domain**

1. Update `AppDomain`/`appDomain` constants. Grep for and update any
   remaining hardcoded strings in code.
1. Set up redirects from old domain to new domain in Go HTTP server.
1. Add the new domain as a custom domain in Cloud Console App Engine
   settings.
1. If using Google Sign In, add the new domain as a valid domain in
   Cloud Console (APIs & Services > OAuth consent screen).
1. Update version number in macOS app for domain change.
