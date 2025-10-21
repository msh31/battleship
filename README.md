# Battleship

A terminal-based implementation of the classic Battleship game. Built with Go using the bubbletea TUI framework and lipgloss for styling.

<img src="https://i.imgur.com/5OrtEh7.png">

## Building

```bash
go build
./battleship
```

## How to Play

The game starts with ship placement. Use arrow keys or WASD to move the cursor, press O to rotate between horizontal and vertical orientation, and hit Space or Enter to place each ship.

Once all five ships are placed, the battle begins. Select a target on the enemy grid and fire. The computer takes its turn after each of your attacks. First player to sink all enemy ships wins.

## Controls

- Arrow keys or WASD: move cursor
- O: toggle ship orientation (placement phase)
- Space/Enter: place ship or fire
- H: show/hide help
- R: restart game
- Q: quit

## Ships

- Carrier: 5 spaces
- Battleship: 4 spaces
- Cruiser: 3 spaces
- Submarine: 3 spaces
- Destroyer: 2 spaces

## Dependencies

- [bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
