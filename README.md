# 🏋️‍♂️ Gymnote - Telegram + VK Bot

GymNote is a workout tracking bot available for **Telegram** and **VK** (ВКонтакте). Designed to help users track their workout progress easily with simple commands and an intuitive interface. Log exercises, track progress over time, and stay consistent with your fitness goals.

> 🆕 **VK Support Added!** Now you can use Gymnote in VK communities. Perfect for users in Russia where Telegram is blocked.

## Features 🏋️‍♂️

- **Workout Logging**: Easily record your exercises with sets, reps, and weights.
- **Progress Tracking**: View statistics on past workouts.
- **Exercise History**: Retrieve past logs to analyze your improvements.
- **User-Friendly Commands**: Simple and efficient command structure.

## Preview

![Preview](/assets/screenshots/preview.gif)

## Tech Stack ⚙️

- **Language**: Go (Golang)
- **Database**: ~~Clickhouse~~ Mongo
- **Cache**: Redis
- **Messaging APIs**: 
  - Telegram Bot API
  - VK API (Long Poll)
- **Hosting**: Deployed on VPS

## Getting Started 🚀

### Prerequisites

Ensure you have the following installed on your system:

- Go (1.23.5)
- Docker & Docker Compose
- Make

### Setup Instructions

1. **Environment Configuration**

```bash
cp .env.example .env
```

2. **Start Services with Docker Compose**

```bash
make docker-up
```

3. **Run migrations** (deprecated)

```bash
make migrate-up
```

4. **Run the Bot**

```bash
make run
```

## Commands 📜

Available in both Telegram and VK:

- **/start** - Start the bot
- **/help** - Show help
- **/start_training** - Start a new training session
- **/upload_training** - Upload a new training session
- **/get_trainings** - View training history
- **/get_exercise_progression** - View weight progression for an exercise
- **/get_exercise_history** - View exercise history
- **/create_exercise** - Create a new exercise (admin only)
- **/clear_training** - Reset the current training session
- **/one_rm** - Calculate one-rep max and percentages

## In action 🚀

### Start Training

![Muscle Groups Screen](/assets/screenshots/start.png)
Kick off your training session by choosing a muscle group. Whether it's chest, back, legs, or arms, GymNote guides you every step of the way.

### Choose Your Exercise

![Exercise Screen](/assets/screenshots/exercise.png)
Browse through a curated list of exercises tailored to your selected muscle group. From bench presses to squats, find the right move for your workout.

### Log Your Sets

![Set Screen](/assets/screenshots/set.png)
Enter your weight and reps for each set. GymNote also shows your exercise history, so you can easily pick the right weight and push your limits.

### Finish Strong

![Finish Screen](/assets/screenshots/finish.png)
At the end of your session, get a detailed summary of your workout. See how many exercises you completed, the total volume lifted, and more.

### Track Your Progression

![Progression Screen](/assets/screenshots/progression.jpg)
Monitor your progress over time with detailed charts. GymNote helps you stay motivated by showing how far you've come in each exercise.

### Training History

![Trainings history](/assets/screenshots/history.png)
Access your complete training history. Review past workouts, analyze your performance, and plan your next session with confidence.

## Deployment ⚙️

The bot is deployed on a **VPS** using Docker and managed via **systemd** for uptime reliability.

### Multi-Platform Support

- **Telegram & VK**: Both bots can run simultaneously
- **Shared Database**: All data is stored in a common MongoDB instance
- **Independent Operation**: Each bot can be enabled/disabled independently via environment variables

This architecture allows easy addition of new platforms (e.g., Discord, Slack) without changing business logic.

## Contributing 🤝

Feel free to open issues, submit pull requests, and improve GymNote together! If you like the project, give it a ⭐ on GitHub!

---

Stay fit and keep logging your progress! 💪🔥
