# 🗣️ Yaketty

**Yaketty** is a Go CLI tool for orchestrating conversations between AI personas using [Ollama](https://ollama.com). Create dialogues between historical figures, celebrities, fictional characters, or entirely original personas in custom scenarios.

## ✨ Features

- 🎭 **Rich Persona Library** - 44 embedded personas including Einstein, Shakespeare, Robin Williams, George Carlin, and many more
- 📜 **Flexible Scenarios** - 13 embedded scenarios from rap battles to philosophical debates, comedy sketches to educational exchanges
- 📦 **Single Binary Portability** - All personas and scenarios embedded at build time
- 🔄 **Local Override Support** - Customize any persona/scenario by placing files in local directories
- 🔧 **Highly Configurable** - YAML-based configuration with command-line overrides
- 🎨 **Creative Combinations** - Mix and match any personas with any scenarios
- 🧠 **Smart Dialogue** - Enhanced prompts ensure character consistency and engaging conversations
- 📚 **Extensible** - Easy to add new personas and scenarios

## 🚀 Quick Start

### Prerequisites

- [Ollama](https://ollama.com) installed and running
- Go 1.22.5+ for building from source

### Installation

```bash
git clone https://github.com/isometry/yaketty.git
cd yaketty
go build -o yaketty .

# List all embedded personas and scenarios
./yaketty list-personas
./yaketty list-scenarios
```

### Your First Dialogue

```bash
# Einstein and Feynman discuss quantum mechanics
./yaketty examples/eli5relativity.yaml

# Eminem vs Tupac rap battle
./yaketty examples/rap-battle-classic.yaml

# Shakespeare and Mark Twain philosophical debate
./yaketty scenarios/philosophy-duel.yaml -1 shakespeare -2 twain
```

## 📖 Configuration

Yaketty uses YAML configuration files that define the scenario and personas:

```yaml
# Basic structure
scenario: |
  The scenario description and context for the conversation.
  This sets the stage and defines the rules of engagement.

opening: |
  The first message that starts the dialogue.

persona1:
  name: "Character Name"
  persona: |
    Detailed character description including personality,
    speech patterns, beliefs, and mannerisms.

persona2:
  name: "Other Character"
  persona: |
    Another character description...

# Optional: AI model parameters
options:
  temperature: 0.8      # Creativity level (0.0-2.0)
  num_ctx: 8192        # Context window size
  top_p: 0.9           # Nucleus sampling
```

## 🎭 Personas

Yaketty includes a rich library of pre-built personas:

### Historical Figures
- **einstein** - Albert Einstein, theoretical physicist and pacifist
- **darwin** - Charles Darwin, methodical naturalist
- **shakespeare** - William Shakespeare, the Bard of Avon
- **voltaire** - Voltaire, Enlightenment philosopher

### Modern Icons
- **obama** - Barack Obama, 44th President
- **jobs** - Steve Jobs, Apple co-founder and visionary
- **attenborough** - David Attenborough, nature documentarian
- **sagan** - Carl Sagan, astronomer and science communicator

### Comedy & Entertainment
- **carlin** - George Carlin, legendary comedian and social critic
- **williams** - Robin Williams, rapid-fire improvisation master
- **elvis** - Elvis Presley, the King of Rock'n'Roll

### Scientists & Thinkers
- **feynman** - Richard Feynman, playful Nobel physicist
- **twain** - Mark Twain, American humorist and satirist

*[View all personas →](personas/)*

## 📚 Scenarios

Pre-built scenarios provide context and structure:

- **🏛️ debate** - Presidential debate format
- **🎤 rap** - Underground rap battle
- **🎭 sketch** - Improvisational comedy
- **🏰 dnd** - Dungeons & Dragons roleplay
- **⚡ enlightenment** - Philosophical podcast
- **🎬 commentary** - Sports commentary
- **🕰️ time-travel-cafe** - Cross-era conversations
- **🔬 philosophy-duel** - AI consciousness debate

*[View all scenarios →](scenarios/)*

## 🛠️ Advanced Usage

### Command Line Options

```bash
# Override personas
./yaketty config.yaml -1 einstein -2 feynman

# Override scenario
./yaketty config.yaml -s philosophy-duel

# Set custom opening line
./yaketty config.yaml -o "Let's discuss the nature of reality"

# Adjust creativity
./yaketty config.yaml --temperature 1.2

# Multiple system prompts
./yaketty config.yaml -p "Be extra witty" -p "Keep responses under 100 words"
```

### Library System

All personas and scenarios are **embedded in the binary** for portability. They can be used by:

**1. Using Embedded Personas/Scenarios** (default):
```bash
# List what's available
./yaketty list-personas
./yaketty list-scenarios

# Use embedded personas
./yaketty config.yaml -1 einstein -2 feynman
```

**2. Local Filesystem Override**:
```bash
# Create custom persona (overrides embedded version if names match)
cat > personas/custom.yaml << EOF
name: My Custom Persona
persona: Custom personality description
EOF

# Use it - loads from local file instead of embedded
./yaketty config.yaml -1 custom
```

**3. Inline in Config**:
```yaml
persona1:
  name: "Inline Persona"
  persona: "Defined directly in config..."
```

**Loading Priority**: Local files → Embedded files → Inline config

### Model Configuration

```yaml
persona1:
  model: llama3        # Specific model for this persona
  options:
    temperature: 0.9     # Per-persona settings
    top_k: 40

persona2:
  model: mistral         # Different model for contrast
  options:
    temperature: 0.7
```

## 💡 Example Combinations

**Educational Dialogues:**
```bash
# Science education
./yaketty time-travel-cafe.yaml -1 einstein -2 feynman

# Historical perspectives
./yaketty enlightenment.yaml -1 voltaire -2 diderot

# Literary analysis
./yaketty philosophy-duel.yaml -1 shakespeare -2 twain
```

**Entertainment:**
```bash
# Comedy gold
./yaketty sketch.yaml -1 carlin -2 williams

# Musical legends
./yaketty rap.yaml -1 elvis -2 eminem

# Presidential perspectives
./yaketty debate.yaml -1 obama -2 biden
```

**Creative Experiments:**
```bash
# Unlikely pairings
./yaketty cooking-disaster.yaml -1 shakespeare -2 jobs

# Cross-domain insights
./yaketty museum-heist.yaml -1 attenborough -2 carlin
```

## 🤝 Contributing

### Adding Personas

Create a new YAML file in `personas/`:

```yaml
name: Your Character
persona: |
  Detailed character description including:
  - Background and expertise
  - Speech patterns and catchphrases
  - Core beliefs and values
  - Personality quirks and mannerisms
  - Communication style
```

### Adding Scenarios

Create a new YAML file in `scenarios/`:

```yaml
scenario: |
  Context and rules for the dialogue.
  What's the setting? What are the goals?

opening: |
  The first message to start the conversation.

# Optional: specific roles
roles:
  - Role for persona1
  - Role for persona2
```

### Guidelines

- **Rich Detail**: Include enough personality details for distinctive voices
- **Clear Context**: Scenarios should provide clear direction and constraints
- **Engaging Openings**: Start with something that immediately draws both personas in
- **Educational Value**: Consider what users might learn from the interaction

## 🎯 Use Cases

- **🎓 Education**: Historical figures explaining concepts to modern students
- **🎪 Entertainment**: Unlikely celebrity conversations and debates
- **💭 Exploration**: Philosophical discussions between great thinkers
- **📝 Writing**: Character development and dialogue practice
- **🧑‍🏫 Teaching**: Interactive lessons with historical personalities
- **🎨 Creativity**: Experimental combinations for artistic inspiration

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details

---

*"The best conversations happen when great minds meet." - Yaketty*

**⭐ Star this repo if you enjoy watching Einstein debate Shakespeare!**
