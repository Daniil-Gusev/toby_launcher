# Toby Launcher

Toby Launcher is a cross-platform command-line interface (CLI) application designed to simplify the use of the [Toby Accessibility Mod for Doom](https://github.com/Alando1-doom/Toby-Accessibility-Mod-for-Doom), making the classic game *Doom* accessible to visually impaired players through text-to-speech (TTS) integration. By providing a unified interface to launch *Doom* games with the accessibility mod, Toby Launcher eliminates the need to manage individual batch files, offering a streamlined experience for players on Windows, macOS, and Linux (x86_64 and arm64).

## Features

- **Cross-Platform Support**: Runs seamlessly on Windows, macOS, and Linux (x86_64 and arm64).
- **Text-to-Speech Integration**: Supports multiple TTS engines to provide audio feedback for visually impaired players:
  - **macOS**: Native NSSpeechSynthesizer.
  - **Windows**: SAPI (Speech API) and NVDA Controller Client for integration with the NVDA screen reader.
  - **Cross-Platform**: eSpeak and eSpeak-NG for all supported platforms.
- **Modular Architecture**: Easily extensible to support additional TTS engines by implementing new packages in `toby_launcher/speech_engines` and registering them in `toby_launcher/speech_engines/engines.go`.
- **Game Management**: Configures and launches *Doom* games using GZDoom with customizable settings stored in a JSON configuration file.
- **Cross-Platform Installer**: Supports both system-wide and portable installations, with an uninstall option for system installations.
- **CLI Interface**: Intuitive menu-driven interface with commands like `help`, `quit`, and `version` for easy navigation.
- **Native Integration**:
  - **macOS**: Packaged as a `.app` bundle for native launching.
  - **Linux**: Creates `.desktop` shortcuts for quick access.
  - **Windows**: Generates `.lnk` shortcuts for desktop integration.

## Installation

### Prerequisites

- **7zip**: Required for the build process to handle archived game data. Install it on your system before building.
- **GZDoom**: The launcher uses [GZDoom](https://github.com/ZDoom/gzdoom) to run *Doom* games. Ensure GZDoom binaries are available or use the provided `download_gzdoom.sh` script to fetch them.
- **Game Files**: You need the necessary `.wad`, `.pk3`, or other game files for *Doom* and the Toby Accessibility Mod.

### Building the Launcher

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/Daniil-Gusev/toby_launcher
   cd toby_launcher
   ```

2. **Prepare Game Files**:
   - Place all required game files (`.wad`, `.pk3`, etc.) in the `resources/data/files` directory.
   - Configure available games in `resources/data/games.json` with the following structure:
     ```json
     {
       "game_name": {
         "description": "Brief description of the game",
         "config": "path/to/config.ini",
         "iwad": "path/to/main.wad",
         "files": ["path/to/additional_file1.pk3", "path/to/additional_file2.wad"],
         "params": ["-param1", "-param2 value"]
       }
     }
     ```

3. **Add Libraries (Windows Only)**:
   - Place required libraries, such as `nvdaControllerClient.dll`, in `resources/lib/<platform_architecture>` (e.g., `resources/lib/windows_amd64`).

4. **Download GZDoom Binaries**:
   - Run the provided script to download GZDoom binaries for all supported platforms:
     ```bash
     ./download_gzdoom.sh
     ```

5. **Build the Launcher**:
   - Execute the build script to create installers for all platforms:
     ```bash
     ./build_installable_release.sh
     ```
   - This script requires `7zip` and generates installers in the project directory.

### Installing the Launcher

The installer supports two modes: **system-wide** and **portable**.

#### System-Wide Installation
- **Windows**: Run the installer with administrator privileges to install to `C:\ProgramData\<AppName>`. A desktop shortcut (`.lnk`) is created.
- **macOS**: Installs to `/Applications/<AppName>.app` with a native `.app` bundle.
- **Linux**: Installs to `/usr/local/bin` with a `.desktop` entry in `/usr/share/applications`.

Run the installer and select option `1` for system installation. Administrator privileges are required.

#### Portable Installation
- Run the installer and select option `2` to install in the current directory. No administrator privileges are needed.
- The application and its data are extracted to `<current_directory>/<AppName>`.

#### Uninstallation
- Select option `3` in the installer to remove a system-wide installation. Administrator privileges are required.

## Usage

1. **Launch the Application**:
   - **macOS**: Double-click the `<AppName>.app` bundle.
   - **Linux**: Use the desktop entry or run the binary from `/usr/local/bin`.
   - **Windows**: Use the desktop shortcut or run the binary directly.
   - **Portable**: Navigate to the installation directory and run the binary.

2. **CLI Interface**:
   - The launcher presents a CLI with a prompt (`> `).
   - Available commands include:
     - `help`: Display available commands.
     - `quit`: Exit the launcher.
     - `version`: Show the launcher version.
     - Custom commands for game selection and management (defined in `core/command.go`).

3. **Game Launch**:
   - Use the CLI to select and start a game configured in `games.json`.
   - The launcher processes game output through the `TextProcessor` (in `core/game/processor.go`), applying rules from `tts_lines.json` to filter and convert text to speech.

4. **Text-to-Speech**:
   - Game output is processed and spoken using the configured TTS engine.
   - Adjust speech rate via the TTS manager if supported by the engine (e.g., NSSpeech, SAPI, eSpeak).

## Project Structure

- **toby_launcher**: Core application logic, including:
  - **core/**: Application logic, CLI interface, and state management.
  - **speech_engines/**: Modular TTS engine implementations (NSSpeech, SAPI, NVDA, eSpeak).
  - **core/game/**: Game management and text processing for accessibility.
- **installer**: Cross-platform installer for system-wide and portable setups.
- **resources/**: Game files, configurations, and platform-specific libraries.
- **download_gzdoom.sh**: Script to fetch GZDoom binaries.

## Adding a New TTS Engine

To add a new TTS engine:
1. Create a new package in `toby_launcher/speech_engines/<engine_name>`.
2. Implement the `SpeechSynthesizer` interface (defined in `core/tts/tts.go`).
3. Register the synthesizer in `speech_engines/engines.go` using `tts.RegisterSynthesizer`.
4. Rebuild the launcher using `build_installable_release.sh`.

Example:
```go
package newengine

import (
    "toby_launcher/core/tts"
)

func init() {
    tts.RegisterSynthesizer(NewSynthesizer, priority)
}

type Synthesizer struct {
    tts.BaseSynthesizer
}

func NewSynthesizer() (tts.SpeechSynthesizer, error) {
    // Implementation
}
```

## Dependencies

- [github.com/chzyer/readline](https://pkg.go.dev/github.com/chzyer/readline): For the CLI interface.
- [github.com/go-ole/go-ole](https://pkg.go.dev/github.com/go-ole/go-ole): For Windows SAPI integration.
- [github.com/bodgit/sevenzip](https://pkg.go.dev/github.com/bodgit/sevenzip): For handling archived game data in the installer.
- [golang.org/x/sys](https://pkg.go.dev/golang.org/x/sys): For system-level operations.

## Acknowledgments

- **[GZDoom Developers](https://github.com/ZDoom/gzdoom)**: For providing the open-source *Doom* engine used by the launcher.
- **[Toby Accessibility Mod Developers](https://github.com/Alando1-doom/Toby-Accessibility-Mod-for-Doom)**: For creating the accessibility mod that enables *Doom* for visually impaired players.
- **[readline](https://pkg.go.dev/github.com/chzyer/readline)**: For powering the CLI interface.
- **[sevenzip](https://pkg.go.dev/github.com/bodgit/sevenzip)**: For enabling cross-platform archive extraction in the installer.

## Contributing

Contributions are welcome! To contribute:
1. Fork the repository: [https://github.com/Daniil-Gusev/toby_launcher](https://github.com/Daniil-Gusev/toby_launcher).
2. Create a feature branch (`git checkout -b feature/your-feature`).
3. Commit your changes (`git commit -m "Add your feature"`).
4. Push to the branch (`git push origin feature/your-feature`).
5. Open a pull request.

Please ensure your code follows the project's coding standards and includes appropriate tests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.