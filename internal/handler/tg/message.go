package tg

const (
	// commands
	startCommand                  = "start"
	helpCommand                   = "help"
	startTrainingCommand          = "start_training"
	clearTrainingCommand          = "clear_training"
	statsCommand                  = "stats"
	createExerciseCommand         = "create_exercise"
	uploadTrainingCommand         = "upload_training"
	getTrainingsCommand           = "get_trainings"
	getExerciseProgressionCommand = "get_exercise_progression"

	// callbacks
	musclePrefix                      = "muscle:"
	exercisePrefix                    = "exercise:"
	finishTrainingPrefix              = "finish_training:"
	startNewExercisePrefix            = "start_new_exercise:"
	startGetExerciseProgressionPrefix = "start_progression:"

	nextDirection = "next"
	prevDirection = "prev"
)

const (
	startText                             = "Я бот для ведения дневника тренировок. Используй команду /help, чтобы узнать доступные команды."
	helpText                              = "📋 Список команд:\n/start - Запустить бота\n/help - Показать справку\n/start_training - Начать новую тренировку\n/upload_training - Загрузить новую тренировку\n/get_trainings - Посмотреть историю тренировок\n/get_exercise_progression - Посмотреть прогрессию весов по упражнению\n/create_exercise - Создать новое упражнение\n/clear_training - Сбросить текущую тренировку\n\nНажимай команды и следуй подсказкам, чтобы вести тренировочный дневник!"
	clearTrainingDoneText                 = "✅ Текущая тренировка успешно удалена!"
	donateAuthorText                      = "\nPS: не забудь подкинуть деньжат %s"
	startTrainingText                     = "🏋️ *Новая тренировка началась!* Выбери мышечную группу:"
	muscleGroupDoneText                   = "✅ Выбрано: *%s*\nТеперь выбери упражнение:"
	muscleGroupSelectText                 = "🏋️ Выбери мышечную группу для нового упражнения:"
	startProgressionMuscleGroupSelectText = "В статистике учитываются тренировки за последний год.\n🏋️ Выбери мышечную группу:"
	exerciseText                          = "✅ Отлично! Вы выбрали упражнение.\nВведите вес и количество повторений через запятую (например: 50.5,12):"
	setText                               = "✅ Подход сохранён! Введите данные нового подхода, либо выберите действие:"
	exerciseCreatedText                   = "Упражнение \"%s\" добавлено в группу \"%s\""
	startNewExerciseText                  = "➕ Начать новое упражнение"
	finishTrainingText                    = "🏁 Завершить тренировку"
	finishText                            = "🏁 Тренировка завершена!\n• Упражнений: %d\n• Подходов: %d\n• Общий вес (кг): %.2f"
	notFoundTrainingsText                 = "🏋️‍♂️ Тренировок пока нет... Но каждый путь начинается с первого шага! Давай, жги, и пусть следующий запрос покажет твои крутые результаты! 🔥"
	startCreateExerciseText               = "Введите название упражнения, группу мышц и оборудование через пробел:\n\nФормат: <название> <группа мышц> <оборудование>"
	startGetTrainingsText                 = "📅 Введите период поиска тренировок в формате: ГГГГ-ММ-ДД ГГГГ-ММ-ДД (например, 2024-12-31 2025-01-22).\nЕсли не укажете даты — покажем тренировки за последние 14 дней. 🔍"
	startUploadTrainingText               = "Введите всю тренировку в формате:\n<год-месяц-число> (опционально)\n<номер упражнения>. <название упражнения> - <вес>,<кол-во повторений> (заметка по подходу); <вес>,<кол-во повторений> (заметка по подходу)\n\nПример:\n2025-01-31\n1. Бабочка - 82,7 (тяжело); 72,8 (тяжело); 54.5,12 (тяжело)\n2. Жим гантелей лежа - 25,10 (нормально); 25,10 (нормально)"
	paginationNextText                    = "Вперед ➡️"
	paginationPrevText                    = "⬅️ Назад"
	loadingProgressionText                = "⏳ График уже строится, ожидайте"

	emptyExerciseNameText             = "Пустое имя упражнения"
	unknownMuscleGroupText            = "Неизвестная мышечная группа!\nДоступные: %v"
	exerciseWithNameAlreadyExistsText = "Упражнение \"%s\" уже существует!"
	unknownCommandText                = "Неизвестная команда. Используй /help для справки"

	// error messages
	errStartTraining     = "❌ Ошибка при запуске тренировки: %v"
	errNoTraining        = "❌ Нет активной тренировки"
	errExerciseLoad      = "❌ Ошибка загрузки упражнений"
	errClearTraining     = "❌ Ошибка сброса тренировки"
	errUploadTraining    = "❌ Ошибка загрузки тренировки"
	errGetTrainings      = "❌ Ошибка поиска тренировок"
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
