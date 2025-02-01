package tg

const (
	// commands
	startCommand          = "start"
	helpCommand           = "help"
	startTrainingCommand  = "start_training"
	clearTrainingCommand  = "clear_training"
	statsCommand          = "stats"
	createExerciseCommand = "create_exercise"

	// callbacks
	musclePrefix           = "muscle:"
	exercisePrefix         = "exercise:"
	finishTrainingPrefix   = "finish_training:"
	startNewExercisePrefix = "start_new_exercise:"

	nextDirection = "next"
	prevDirection = "prev"
)

const (
	// success messages
	startText             = "Я бот для ведения дневника тренировок. Используй команду /help, чтобы узнать доступные команды."
	helpText              = "📋 Список команд:\n/start - Запустить бота\n/help - Показать справку\n/start_training - Начать новую тренировку\n/create_exercise - Создать новое упражнение\n/clear_training - Сбросить текущую тренировку\n\nНажимай команды и следуй подсказкам, чтобы вести тренировочный дневник!"
	clearTrainingDoneText = "✅ Текущая тренировка успешно удалена!"
	donateAuthorText      = "\nPS: не забудь подкинуть деньжат %s"
	startTrainingText     = "🏋️ *Новая тренировка началась!* Выбери мышечную группу:"
	muscleGroupDoneText   = "✅ Выбрано: *%s*\nТеперь выбери упражнение:"
	muscleGroupSelectText = "🏋️ Выбери мышечную группу для нового упражнения:"
	exerciseText          = "✅ Отлично! Вы выбрали упражнение.\nВведите вес и количество повторений через запятую (например: 50.5,12):"
	setText               = "✅ Подход сохранён! Введите данные нового подхода, либо выберите действие:"
	exerciseCreatedText   = "Упражнение \"%s\" добавлено в группу \"%s\""
	startNewExerciseText  = "➕ Начать новое упражнение"
	finishTrainingText    = "🏁 Завершить тренировку"
	finishText            = "🏁 Тренировка завершена!\n• Упражнений: %d\n• Подходов: %d\n• Общий вес (кг): %.2f"
	startCreateExercise   = "Введите название упражнения, группу мышц и оборудование через пробел:\n\nФормат: <название> <группа мышц> <оборудование>"
	paginationNextText    = "Вперед ➡️"
	paginationPrevText    = "⬅️ Назад"

	emptyExerciseNameText             = "Пустое имя упражнения"
	unknownMuscleGroupText            = "Неизвестная мышечная группа!\nДоступные: %v"
	exerciseWithNameAlreadyExistsText = "Упражнение \"%s\" уже существует!"
	unknownCommandText                = "Неизвестная команда. Используй /help для справки"

	// error messages
	errStartTraining     = "❌ Ошибка при запуске тренировки: %v"
	errNoTraining        = "❌ Нет активной тренировки"
	errExerciseLoad      = "❌ Ошибка загрузки упражнений"
	errClearTraining     = "❌ Ошибка сброса тренировки"
	errAddExercise       = "❌ Ошибка при добавлении упражнения: %v"
	errInvalidFormat     = "❌ Неверный формат. Введите вес и повторения через запятую (например: 50.5,12)"
	errParseData         = "❌ Ошибка при разборе данных. Проверьте формат и попробуйте снова."
	errGeneral           = "❌ Ошибка: %v"
	errInvalidExerciseID = "❌ Ошибка: неверный формат ID упражнения."
	errCreateExercise    = "❌ Ошибка при добавлении упражнения"
)

var (
	muscleGroupsWithSmiles = []string{"💪 Спина", "🏋️ Грудь", "🦵 Ноги", "💪 Руки"}
	availableMuscleGroups  = []string{"Спина", "Грудь", "Ноги", "Руки"}
)
