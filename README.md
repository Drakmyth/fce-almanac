# fce-almanac

An encyclopedia of information about [FortressCraft Evolved](https://store.steampowered.com/app/254200/FortressCraft_Evolved/) providing information critical to playing efficiently in an easily digestible format.

### Built With

* [![Htmx][htmx-shield]][htmx-url]
* [![Golang][golang-shield]][golang-url]
* [![TailwindCSS][tailwindcss-shield]][tailwindcss-url]

## Getting Started

### Prerequisites

This site only requires the Golang compiler, tooling, and libraries. These can be downloaded from the [Golang homepage](https://go.dev/) or installed via a package manager:
* apt
  ```sh
  > apt install golang-go
  ```
* Chocolatey
  ```pwsh
  > choco install -y golang
  ```
* Homebrew
  ```sh
  > brew install go
  ```
* winget
  ```pwsh
  > winget install GoLang.Go
  ```

### Installation & Execution

There is currently no additional installation necessary. Simply clone the repo and start the server.

1. Clone the repo
   ```sh
   > git clone git@github.com:Drakmyth/fce-almanac.git
   ```
1. Start the server
   ```sh
   > go run server.go
   ```

## Contributing

While the data provided by this site is sourced from game files and is not provided here, the server software itself is open source! If you would like to contribute please feel welcome to do so! If you have a suggestion, you can fork the repo and create a pull request. Alternatively you can open an issue with the `enhancement` tag.

## License

Distributed under the MIT License. See [LICENSE.md](./LICENSE.md) for more information.


<!-- Reference Links -->
[htmx-url]: https://htmx.org
[htmx-shield]: https://img.shields.io/badge/htmx-4470d2?style=for-the-badge
[golang-url]: https://go.dev
[golang-shield]: https://img.shields.io/badge/golang-09657c?style=for-the-badge&logo=go&logoColor=79d2fa
[tailwindcss-url]: https://tailwindcss.com/
[tailwindcss-shield]: https://img.shields.io/badge/Tailwind%20CSS-0b111f?style=for-the-badge&logo=tailwindcss&logoColor=26bcf5

