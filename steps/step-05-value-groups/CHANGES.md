# step-05-value-groups

Introduces grouped route registration.

- What changed: handlers are provided into `group:"routes"`; mux consumes the grouped slice.
- Why it matters: adding a route no longer requires editing central registration code.
- Verify: both `POST /shorten` and `GET /{slug}` are registered and functional.
