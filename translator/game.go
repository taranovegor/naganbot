package translator

var GameTranslations = translations{
	"ru": {
		"something went wrong": {message: SimpleMessage("Наган заклинило. Мы уже вызвали оружейника.")},
		"game creation": {
			oneOf: []oneOf{
				{message: SimpleMessage("Мы играем чтобы пощекотать нервы, хочешь тоже?")},
				{message: SimpleMessage("Я вижу ты не из трусливых, что насчёт остальных?")},
				{message: SimpleMessage("Вы опять за старое? Так уж и быть.")},
				{message: SimpleMessage("Ну что, рискнёте сыграть в игру с летальным исходом?")},
				{message: SimpleMessage("Смерть или слава? Давайте узнаем!")},
				{message: SimpleMessage("Кто-то сегодня не вернётся домой...")},
				{message: SimpleMessage("Пистолет заряжен, очередь за вами.")},
				{message: SimpleMessage("Русская рулетка — лучший способ проверить удачу.")},
				{message: SimpleMessage("Готовы испытать судьбу?")},
				{message: SimpleMessage("Один выстрел — и ты в истории.")},
				{message: SimpleMessage("Судьба любит смелых. Или нет?")},
				{message: SimpleMessage("На кону жизнь. Ну, почти.")},
			},
		},
		"joining the game": {
			oneOf: []oneOf{
				{message: SimpleMessage("Могу поспорить, что ты - труп.")},
				{message: SimpleMessage("Ты считаешь себя бессмертным?")},
				{message: SimpleMessage("Ну хорошо, присаживайся, игра скоро начнётся.")},
				{message: SimpleMessage("Решил напоследок рискнуть?")},
				{message: SimpleMessage("1 к 6 - вот твой шанс выживания в русской рулетке.")},
				{message: SimpleMessage("Здесь произойдёт охуенное убийство.")},
				{message: SimpleMessage("Ты играл когда-нибудь до этого в рулетку?")},
				{message: SimpleMessage("Добро пожаловать в клуб самоубийц.")},
				{message: SimpleMessage("Шанс выжить? Примерно как у снега в аду.")},
				{message: SimpleMessage("Ты либо герой, либо покойник.")},
				{message: SimpleMessage("Один патрон, шесть камер. Удачи.")},
				{message: SimpleMessage("Ты уверен, что хочешь это сделать?")},
				{message: SimpleMessage("Сейчас узнаем, насколько ты везучий.")},
				{message: SimpleMessage("Ты либо выйдешь победителем, либо станешь уроком для остальных.")},
				{message: SimpleMessage("Ну что, готов к своему последнему клику?")},
				{message: SimpleMessage("Судьба улыбается... или нет?")},
			},
		},
		"play the game": {
			oneOf: []oneOf{
				{
					allOf: []oneOf{
						{message: SimpleMessage("Решили рискнуть? Тогда начнём.")},
						{message: SimpleMessage("Стопка, щелчок, выстрел...")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Хорошо, раз все в сборе, начнём.")},
						{message: SimpleMessage("Пистолет передаётся следующему...")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Игра началась. Пусть удача будет на вашей стороне.")},
						{message: SimpleMessage("Колесо фортуны крутится...")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Ну что, паника или спокойствие?")},
						{message: SimpleMessage("Пистолет уже в игре...")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Время узнать, кто сегодня проиграл.")},
						{message: SimpleMessage("Палец на спусковом курке...")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Готовы к самому страшному клику в вашей жизни?")},
						{message: SimpleMessage("Курок нажат...")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Сейчас всё решит один выстрел.")},
						{message: SimpleMessage("Тишина перед бурей...")},
					},
				}},
		},
		"gunslinger killed": {
			oneOf: []oneOf{
				{message: SimpleMessage("Раздался выстрел — и %gunslinger больше с нами.")},
				{message: SimpleMessage("Барабан провернулся, курок щёлкнул... %gunslinger выбывает.")},
				{message: SimpleMessage("Пуля нашла свою жертву — это был %gunslinger.")},
				{message: SimpleMessage("Кровь, порох и тишина... %gunslinger проиграл.")},
				{message: SimpleMessage("%gunslinger решил, что жил слишком долго.")},
				{message: SimpleMessage("В этот раз не повезло %gunslinger.")},
				{message: SimpleMessage("На этот раз смерть выбрала %gunslinger.")},
				{message: SimpleMessage("Роковой выстрел для %gunslinger.")},
				{message: SimpleMessage("%gunslinger больше не ответит.")},
				{message: SimpleMessage("Прощай, %gunslinger, мы тебя запомним...")},
				{message: SimpleMessage("Игра окончена для %gunslinger.")},
				{message: SimpleMessage("Момент истины для %gunslinger оказался последним.")},
				{message: SimpleMessage("%gunslinger получил билет в один конец.")},
			},
		},
		"killed by atomic bullet": {
			oneOf: []oneOf{
				{message: SimpleMessage("Щёлк — и внезапно комната превратилась в эпицентр ядерного гриба. В револьере оказалась атомная пуля.")},
				{message: SimpleMessage("...от группы смельчаков остался только радиоактивный след. Видимо, кто-то подсунул в барабан не ту пулю.")},
				{message: SimpleMessage("Ккомната наполнилась ярким светом... Все игроки мгновенно испарились. Кто-то явно жульничал с боеприпасами.")},
			},
		},
		"joined the game":       {message: SimpleMessage("Участники игры %date.")},
		"game join list item":   {message: SimpleMessage("%num. %gunslinger.")},
		"owner of the game":     {message: SimpleMessage("владелец.")},
		"shot in game":          {message: SimpleMessage("выбыл.")},
		"top is not determined": {message: SimpleMessage("Пока что мы не можем составить топ игроков.")},
		"top players by games":  {message: SimpleMessage("Топ-%number выбывших игроков:.")},
		"top game player": {message: PluralMessage{
			one:  "%i. %user - %times раз",
			few:  "%i. %user - %times раза",
			many: "%i. %user - %times раз",
		}},
		"available only in chat": {message: SimpleMessage("Русская рулетка игра не для одного.")},
		"player already in game": {
			oneOf: []oneOf{
				{message: SimpleMessage("Ты уже в игре. Смерть тебя уже запомнила.")},
				{message: SimpleMessage("Жаждешь смерти? Подожди своей очереди.")},
				{message: SimpleMessage("Ты уже в списке на расстрел. Терпение.")},
				{message: SimpleMessage("Ты уже в деле. Расслабься и жди выстрела.")},
				{message: SimpleMessage("Ты уже участвуешь. Или передумал?")},
				{message: SimpleMessage("Смерть любит нетерпеливых.")},
				{message: SimpleMessage("Ты уже в игре, не торопи смерть.")},
				{message: SimpleMessage("Один раз в игре - достаточно для проверки удачи.")},
				{message: SimpleMessage("Ты уже в списке потенциальных покойников.")},
			},
		},
		"player is not kicked": {
			oneOf: []oneOf{
				{message: SimpleMessage("Револьвер заклинило. Видимо, у кого-то иммунитет.")},
				{message: SimpleMessage("Порох отсырел — приговор не был приведён в исполнение.")},
				{message: SimpleMessage("Смерть сегодня капризничает. Кое-кто временно неуязвим.")},
				{message: SimpleMessage("Казнь откладывается — технические неполадки.")},
				{message: SimpleMessage("Кто-то подложил холостой патрон.")},
			},
		},
		"wait for game timeout": {
			oneOf: []oneOf{
				{message: SimpleMessage("Дайте барабану остыть после последнего выстрела.")},
				{message: SimpleMessage("Смерть тоже любит делать перерывы, не спешите.")},
				{message: SimpleMessage("Подождите, пока духи прошлых игроков разойдутся.")},
				{message: SimpleMessage("Дайте нам время замести следы от прошлой игры.")},
				{message: SimpleMessage("Мы наделали слишком много шума, пока нам нужно залечь на дно.")},
				{message: SimpleMessage("Смерть пока занята другими.")},
			},
		},
		"active game not found":                  {message: SimpleMessage("Ещё никто не начал игру.")},
		"user game statistics":                   {message: SimpleMessage("Вы участвовали в %games игре(-ах) и проиграли %shots раз(а).")},
		"available settings below":               {message: SimpleMessage("Персонализируйте игру используя перечисленные ниже опции.")},
		"settings can be changed only by admins": {message: SimpleMessage("Изменять настройки игры могут только администраторы чата.")},
		"4 shot revolver":                        {message: SimpleMessage("Colt Cloverleaf - 4 игрока")},
		"6 shot revolver":                        {message: SimpleMessage("Colt Python - 6 игроков")},
		"7 shot revolver":                        {message: SimpleMessage("Наган - 7 игроков")},
		"revolver has been replaced": {message: PluralMessage{
			one:  "Теперь в игре может участвовать %i игрок",
			few:  "Теперь в игре могут участвовать %i игрока",
			many: "Теперь в игре могут участвовать %i игроков",
		}},
		"settings will be applied for next games": {message: SimpleMessage("Изменения вступят в силу со следующего раунда")},
	},
	"uk": {
		"something went wrong": {message: SimpleMessage("Щось пішло не так. Мабуть, наган знову клинить.")},
		"game creation": {
			oneOf: []oneOf{
				{message: SimpleMessage("Граємо в гру зі смертельним кінцем. Приєднаєшся?")},
				{message: SimpleMessage("Бачу, ти сміливий. А як щодо інших?")},
				{message: SimpleMessage("Знову за старе? Ну що ж, поїхали.")},
				{message: SimpleMessage("Ну що, ризикнете зіграти в гру з летальним результатом?")},
				{message: SimpleMessage("Смерть чи слава? Давайте дізнаємось!")},
				{message: SimpleMessage("Хтось сьогодні не повернеться додому...")},
				{message: SimpleMessage("Пістолет заряджений, черга за вами.")},
				{message: SimpleMessage("Рулетка - найкращий спосіб перевірити вдачу.")},
				{message: SimpleMessage("Готові випробувати долю?")},
				{message: SimpleMessage("Один постріл - і ти в історії.")},
				{message: SimpleMessage("Доля любить сміливців. Чи ні?")},
				{message: SimpleMessage("На кону життя. Ну, майже.")},
			},
		},
		"joining the game": {
			oneOf: []oneOf{
				{message: SimpleMessage("Можу поставити, що ти - труп.")},
				{message: SimpleMessage("Ти вважаєш себе безсмертним?")},
				{message: SimpleMessage("Добре, сідай, гра скоро почнеться.")},
				{message: SimpleMessage("Вирішив напослідок ризикнути?")},
				{message: SimpleMessage("1 до 6 - ось твій шанс вижити в рулетці.")},
				{message: SimpleMessage("Тут станеться огидне вбивство.")},
				{message: SimpleMessage("Ти коли-небудь грав у рулетку?")},
				{message: SimpleMessage("Ласкаво просимо до клубу самогубців.")},
				{message: SimpleMessage("Шанс вижити? Приблизно як у снігу в пеклі.")},
				{message: SimpleMessage("Ти або герой, або покійник.")},
				{message: SimpleMessage("Один набій, шість камер. Удачі.")},
				{message: SimpleMessage("Ти впевнений, що хочеш це зробити?")},
				{message: SimpleMessage("Зараз дізнаємось, наскільки ти везучий.")},
				{message: SimpleMessage("Ти або вийдеш переможцем, або станеш уроком для інших.")},
				{message: SimpleMessage("Ну що, готовий до свого останнього кліку?")},
				{message: SimpleMessage("Доля посміхається... чи ні?")},
			},
		},
		"play the game": {
			oneOf: []oneOf{
				{
					allOf: []oneOf{
						{message: SimpleMessage("Вирішили ризикнути? Тоді починаємо.")},
						{message: SimpleMessage("Барабан, клац, постріл...")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Добре, раз всі в зборі, почнемо.")},
						{message: SimpleMessage("Пістолет передається наступному...")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Гра почалась. Нехай удача буде на вашій стороні.")},
						{message: SimpleMessage("Колесо фортуни крутиться...")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Ну що, паніка чи спокій?")},
						{message: SimpleMessage("Пістолет уже в грі...")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Час дізнатись, хто сьогодні програв.")},
						{message: SimpleMessage("Палець на спусковому курку...")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Готові до найстрашнішого кліку у вашому житті?")},
						{message: SimpleMessage("Курок натиснуто...")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Зараз все вирішить один постріл.")},
						{message: SimpleMessage("Тиша перед бурею...")},
					},
				}},
		},
		"gunslinger killed": {
			oneOf: []oneOf{
				{message: SimpleMessage("Лунає постріл — і %gunslinger більше не з нами.")},
				{message: SimpleMessage("Барабан прокрутився, курок клацнув... %gunslinger вибуває.")},
				{message: SimpleMessage("Куля знайшла свою жертву — це був %gunslinger.")},
				{message: SimpleMessage("Кров, порох і тиша... %gunslinger програв.")},
				{message: SimpleMessage("%gunslinger вирішив, що жив занадто довго.")},
				{message: SimpleMessage("Цього разу не пощастило %gunslinger.")},
				{message: SimpleMessage("На цей раз смерть обрала %gunslinger.")},
				{message: SimpleMessage("Фатальний постріл для %gunslinger.")},
				{message: SimpleMessage("%gunslinger більше не відповість.")},
				{message: SimpleMessage("Прощавай, %gunslinger, ми тебе запам'ятаємо...")},
				{message: SimpleMessage("Гра закінчена для %gunslinger.")},
				{message: SimpleMessage("Момент істини для %gunslinger виявився останнім.")},
				{message: SimpleMessage("%gunslinger отримав квиток в один кінець.")},
			},
		},
		"killed by atomic bullet": {
			oneOf: []oneOf{
				{message: SimpleMessage("Клац — і раптом кімната перетворилась на епіцентр ядерного гриба. В револьвері виявилась атомна куля.")},
				{message: SimpleMessage("...від групи сміливців залишився лише радіоактивний слід. Мабуть, хтось підсунув у барабан не той набій.")},
				{message: SimpleMessage("Кімната наповнилась яскравим світлом... Усі гравці миттєво випарувались. Хтось явно шахраював з боєприпасами.")},
			},
		},
		"joined the game":       {message: SimpleMessage("Учасники гри %date.")},
		"game join list item":   {message: SimpleMessage("%num. %gunslinger.")},
		"owner of the game":     {message: SimpleMessage("власник.")},
		"shot in game":          {message: SimpleMessage("вибув.")},
		"top is not determined": {message: SimpleMessage("Поки що ми не можемо скласти топ гравців.")},
		"top players by games":  {message: SimpleMessage("Топ-%number вибулих гравців:.")},
		"top game player": {message: PluralMessage{
			one:  "%i. %user - %times раз",
			few:  "%i. %user - %times рази",
			many: "%i. %user - %times разів",
		}},
		"available only in chat": {message: SimpleMessage("Рулетка - гра не для одного.")},
		"player already in game": {
			oneOf: []oneOf{
				{message: SimpleMessage("Ти вже в грі. Смерть тебе вже запам'ятала.")},
				{message: SimpleMessage("Прагнеш смерті? Почекай своєї черги.")},
				{message: SimpleMessage("Ти вже в списку на розстріл. Терпіння.")},
				{message: SimpleMessage("Ти вже в справі. Розслабся і чекай пострілу.")},
				{message: SimpleMessage("Ти вже береш участь. Чи передумав?")},
				{message: SimpleMessage("Смерть любить нетерплячих.")},
				{message: SimpleMessage("Ти вже в грі, не поспішай смерть.")},
				{message: SimpleMessage("Один раз у грі - достатньо для перевірки удачі.")},
				{message: SimpleMessage("Ти вже в списку потенційних покійників.")},
			},
		},
		"player is not kicked": {
			oneOf: []oneOf{
				{message: SimpleMessage("Револьвер заклинено. Мабуть, у когось імунітет.")},
				{message: SimpleMessage("Порох відмок — вирок не було виконано.")},
				{message: SimpleMessage("Смерть сьогодні капризничає. Дехто тимчасово невразливий.")},
				{message: SimpleMessage("Страту відкладено — технічні неполадки.")},
				{message: SimpleMessage("Хтось підклав холостий набій.")},
			},
		},
		"wait for game timeout": {
			oneOf: []oneOf{
				{message: SimpleMessage("Дайте барабану охолонути після останнього пострілу.")},
				{message: SimpleMessage("Смерть теж любить робити перерви, не поспішайте.")},
				{message: SimpleMessage("Почекайте, поки духи минулих гравців розійдуться.")},
				{message: SimpleMessage("Дайте нам час замести сліди від минулої гри.")},
				{message: SimpleMessage("Ми наробили занадто багато шуму, поки нам потрібно залягти на дно.")},
				{message: SimpleMessage("Смерть поки зайнята іншими.")},
			},
		},
		"active game not found":                  {message: SimpleMessage("Ще ніхто не почав гру.")},
		"user game statistics":                   {message: SimpleMessage("Ви брали участь у %games грі(-ах) і програли %shots раз(и).")},
		"available settings below":               {message: SimpleMessage("Персоналізуйте гру, використовуючи перелічені нижче опції.")},
		"settings can be changed only by admins": {message: SimpleMessage("Змінювати налаштування гри можуть лише адміністратори чату.")},
		"4 shot revolver":                        {message: SimpleMessage("Colt Cloverleaf - 4 гравці")},
		"6 shot revolver":                        {message: SimpleMessage("Colt Python - 6 гравців")},
		"7 shot revolver":                        {message: SimpleMessage("Наган - 7 гравців")},
		"revolver has been replaced": {message: PluralMessage{
			one:  "Тепер у грі може брати участь %i гравець",
			few:  "Тепер у грі можуть брати участь %i гравці",
			many: "Тепер у грі можуть брати участь %i гравців",
		}},
		"settings will be applied for next games": {message: SimpleMessage("Зміни набудуть чинності з наступного раунду")},
	},
}
