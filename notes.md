# FitHub

## Work in Progress

- Exercises page
  - [ ] List of exercises
    - [ ] By Muscle Group
    - [ ] By Exercise type?
    - [ ] By Level?
  - [ ] Search feature

- Exercise page
  - [ ] Name
  - [ ] Primary/Secondary muscle groups
  - [ ] Description / Instructions
  - [ ] Video / embedded video / link to video
    - [ ] need to decide how I want to handle video instructions
  - [ ] Add to exercise button?

- Implementing the use of refresh tokens
  - [ ] Access tokens expiry times shortened
  - [ ] Handlers needing access tokens utilize refresh handler for new access
        tokens
  - [ ] Revoke refresh token on logout

- Styling and theming

## Known Issues

- Exercise Quick Search on workout page formatting for scroll is awkward.

## Helpful Links

## Future Plans

- Looking to transition auth pattern to hs256 or sesssion based
  - Turns out that rs256 doesn't make much sense with my monolith
  - Interested in knowing more about the session based auth
