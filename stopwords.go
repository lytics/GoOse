package goose

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"gopkg.in/fatih/set.v0"
)

var punctuationRegex = regexp.MustCompile(`[^\p{Ll}\p{Lu}\p{Lt}\p{Lo}\p{Nd}\p{Pc}\s]`)

// StopWords implements a simple language detector
type StopWords struct {
	cachedStopWords map[string]*set.Set
}

// NewStopwords returns an instance of a stop words detector
func NewStopwords() StopWords {
	cachedStopWords := make(map[string]*set.Set)
	for lang, stopwords := range sw {
		lines := strings.Split(stopwords, "\n")
		cachedStopWords[lang] = set.New()
		for _, line := range lines {
			if strings.HasPrefix(line, "#") {
				continue
			}
			line = strings.TrimSpace(line)
			cachedStopWords[lang].Add(line)
		}
	}
	return StopWords{
		cachedStopWords: cachedStopWords,
	}
}

/*
func NewStopwords(path string) StopWords {
	cachedStopWords := make(map[string]*set.Set)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err.Error())
	}
	for _, file := range files {
		name := strings.Replace(file.Name(), ".txt", "", -1)
		name = strings.Replace(name, "stopwords-", "", -1)
		name = strings.ToLower(name)

		stops := set.New()
		lines := ReadLinesOfFile(path + "/" + file.Name())
		for _, line := range lines {
			line = strings.Trim(line, " ")
			stops.Add(line)
		}
		cachedStopWords[name] = stops
	}

	return StopWords{
		cachedStopWords: cachedStopWords,
	}
}
*/

func (stop *StopWords) removePunctuation(text string) string {
	return punctuationRegex.ReplaceAllString(text, "")
}

func (stop *StopWords) stopWordsCount(lang string, text string) wordStats {
	if text == "" {
		return wordStats{}
	}
	ws := wordStats{}
	stopWords := set.New()
	text = strings.ToLower(text)
	items := strings.Split(text, " ")
	stops := stop.cachedStopWords[lang]
	count := 0
	if stops != nil {
		for _, item := range items {
			if stops.Has(item) {
				stopWords.Add(item)
				count++
			}
		}
	}

	ws.stopWordCount = stopWords.Size()
	ws.wordCount = len(items)
	ws.stopWords = stopWords

	return ws
}

// SimpleLanguageDetector returns the language code for the text, based on its stop words
func (stop StopWords) SimpleLanguageDetector(text string) string {
	max := 0
	currentLang := "en"

	for k := range sw {
		ws := stop.stopWordsCount(k, text)
		if ws.stopWordCount > max {
			max = ws.stopWordCount
			currentLang = k
		}
	}

	return currentLang
}

// ReadLinesOfFile returns the lines from a file as a slice of strings
func ReadLinesOfFile(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err.Error())
	}
	lines := strings.Split(string(content), "\n")
	return lines
}

var sw = map[string]string{
	"ar": `
فى
في
كل
لم
لن
له
من
هو
هي
قوة
كما
لها
منذ
وقد
ولا
نفسه
لقاء
مقابل
هناك
وقال
وكان
نهاية
وقالت
وكانت
للامم
فيه
كلم
لكن
وفي
وقف
ولم
ومن
وهو
وهي
يوم
فيها
منها
مليار
لوكالة
يكون
يمكن
مليون
حيث
اكد
الا
اما
امس
السابق
التى
التي
اكثر
ايار
ايضا
ثلاثة
الذاتي
الاخيرة
الثاني
الثانية
الذى
الذي
الان
امام
ايام
خلال
حوالى
الذين
الاول
الاولى
بين
ذلك
دون
حول
حين
الف
الى
انه
اول
ضمن
انها
جميع
الماضي
الوقت
المقبل
اليوم
ـ
ف
و
و6
قد
لا
ما
مع
مساء
هذا
واحد
واضاف
واضافت
فان
قبل
قال
كان
لدى
نحو
هذه
وان
واكد
كانت
واوضح
مايو
ب
ا
أ
،
عشر
عدد
عدة
عشرة
عدم
عام
عاما
عن
عند
عندما
على
عليه
عليها
زيارة
سنة
سنوات
تم
ضد
بعد
بعض
اعادة
اعلنت
بسبب
حتى
اذا
احد
اثر
برس
باسم
غدا
شخصا
صباح
اطار
اربعة
اخرى
بان
اجل
غير
بشكل
حاليا
بن
به
ثم
اف
ان
او
اي
بها
صفر
	`,
	"en": `
	a's
able
about
above
according
accordingly
across
actually
after
afterwards
again
against
ain't
all
allow
allows
almost
alone
along
already
also
although
always
am
among
amongst
an
and
another
any
anybody
anyhow
anyone
anything
anyway
anyways
anywhere
apart
appear
appreciate
appropriate
are
aren't
around
as
aside
ask
asking
associated
at
available
away
awfully
be
became
because
become
becomes
becoming
been
before
beforehand
behind
being
believe
below
beside
besides
best
better
between
beyond
both
brief
but
by
c
c'mon
c's
came
campaign
can
can't
cannot
cant
cause
causes
certain
certainly
changes
clearly
co
com
come
comes
concerning
consequently
consider
considering
contain
containing
contains
corresponding
could
couldn't
course
currently
definitely
described
despite
did
didn't
different
do
does
doesn't
doing
don't
done
down
downwards
during
each
edu
eight
either
else
elsewhere
enough
endorsed
entirely
especially
et
etc
even
ever
every
everybody
everyone
everything
everywhere
ex
exactly
example
except
far
few
fifth
first
financial
five
followed
following
follows
for
former
formerly
forth
four
from
further
furthermore
get
gets
getting
given
gives
go
goes
going
gone
got
gotten
greetings
had
hadn't
happens
hardly
has
hasn't
have
haven't
having
he
he's
hello
help
hence
her
here
here's
hereafter
hereby
herein
hereupon
hers
herself
hi
him
himself
his
hither
hopefully
how
howbeit
however
i'd
i'll
i'm
i've
if
ignored
immediate
in
inasmuch
inc
indeed
indicate
indicated
indicates
inner
insofar
instead
into
inward
is
isn't
it
it'd
it'll
it's
its
itself
just
keep
keeps
kept
know
knows
known
last
lately
later
latter
latterly
least
less
lest
let
let's
like
liked
likely
little
look
looking
looks
ltd
mainly
many
may
maybe
me
mean
meanwhile
merely
might
more
moreover
most
mostly
much
must
my
myself
name
namely
nd
near
nearly
necessary
need
needs
neither
never
nevertheless
new
next
nine
no
nobody
non
none
noone
nor
normally
not
nothing
novel
now
nowhere
obviously
of
off
often
oh
ok
okay
old
on
once
one
ones
only
onto
or
other
others
otherwise
ought
our
ours
ourselves
out
outside
over
overall
own
particular
particularly
per
perhaps
placed
please
plus
possible
presumably
probably
provides
quite
quote
quarterly
rather
really
reasonably
regarding
regardless
regards
relatively
respectively
right
said
same
saw
say
saying
says
second
secondly
see
seeing
seem
seemed
seeming
seems
seen
self
selves
sensible
sent
serious
seriously
seven
several
shall
she
should
shouldn't
since
six
so
some
somebody
somehow
someone
something
sometime
sometimes
somewhat
somewhere
soon
sorry
specified
specify
specifying
still
sub
such
sup
sure
t's
take
taken
tell
tends
than
thank
thanks
thanx
that
that's
thats
the
their
theirs
them
themselves
then
thence
there
there's
thereafter
thereby
therefore
therein
theres
thereupon
these
they
they'd
they'll
they're
they've
think
third
this
thorough
thoroughly
those
though
three
through
throughout
thru
thus
to
together
too
took
toward
towards
tried
tries
truly
try
trying
twice
two
under
unfortunately
unless
unlikely
until
unto
up
upon
us
use
used
useful
uses
using
usually
uucp
value
various
very
via
viz
vs
want
wants
was
wasn't
way
we
we'd
we'll
we're
we've
welcome
well
went
were
weren't
what
what's
whatever
when
whence
whenever
where
where's
whereafter
whereas
whereby
wherein
whereupon
wherever
whether
which
while
whither
who
who's
whoever
whole
whom
whose
why
will
willing
wish
with
within
without
won't
wonder
would
would
wouldn't
yes
yet
you
you'd
you'll
you're
you've
your
yours
yourself
yourselves
zero
official
sharply
criticized
`,
	"es": `
de
la
que
el
en
y
a
los
del
se
las
por
un
para
con
no
una
su
al
lo
como
más
pero
sus
le
ya
o
este
sí
porque
esta
entre
cuando
muy
sin
sobre
también
me
hasta
hay
donde
quien
desde
todo
nos
durante
todos
uno
les
ni
contra
otros
ese
eso
ante
ellos
e
esto
mí
antes
algunos
qué
unos
yo
otro
otras
otra
él
tanto
esa
estos
mucho
quienes
nada
muchos
cual
poco
ella
estar
estas
algunas
algo
nosotros
mi
mis
tú
te
ti
tu
tus
ellas
nosotras
vosotros
vosotras
os
mío
mía
míos
mías
tuyo
tuya
tuyos
tuyas
suyo
suya
suyos
suyas
nuestro
nuestra
nuestros
nuestras
vuestro
vuestra
vuestros
vuestras
esos
esas
estoy
estás
está
estamos
estáis
están
esté
estés
estemos
estéis
estén
estaré
estarás
estará
estaremos
estaréis
estarán
estaría
estarías
estaríamos
estaríais
estarían
estaba
estabas
estábamos
estabais
estaban
estuve
estuviste
estuvo
estuvimos
estuvisteis
estuvieron
estuviera
estuvieras
estuviéramos
estuvierais
estuvieran
estuviese
estuvieses
estuviésemos
estuvieseis
estuviesen
estando
estado
estada
estados
estadas
estad
he
has
ha
hemos
habéis
han
haya
hayas
hayamos
hayáis
hayan
habré
habrás
habrá
habremos
habréis
habrán
habría
habrías
habríamos
habríais
habrían
había
habías
habíamos
habíais
habían
hube
hubiste
hubo
hubimos
hubisteis
hubieron
hubiera
hubieras
hubiéramos
hubierais
hubieran
hubiese
hubieses
hubiésemos
hubieseis
hubiesen
habiendo
habido
habida
habidos
habidas

# forms of ser, to be (not including the infinitive):
soy
eres
es
somos
sois
son
sea
seas
seamos
seáis
sean
seré
serás
será
seremos
seréis
serán
sería
serías
seríamos
seríais
serían
era
eras
éramos
erais
eran
fui
fuiste
fue
fuimos
fuisteis
fueron
fuera
fueras
fuéramos
fuerais
fueran
fuese
fueses
fuésemos
fueseis
fuesen
siendo
sido
tengo
tienes
tiene
tenemos
tenéis
tienen
tenga
tengas
tengamos
tengáis
tengan
tendré
tendrás
tendrá
tendremos
tendréis
tendrán
tendría
tendrías
tendríamos
tendríais
tendrían
tenía
tenías
teníamos
teníais
tenían
tuve
tuviste
tuvo
tuvimos
tuvisteis
tuvieron
tuviera
tuvieras
tuviéramos
tuvierais
tuvieran
tuviese
tuvieses
tuviésemos
tuvieseis
tuviesen
teniendo
tenido
tenida
tenidos
tenidas
tened
`,
	"fr": `
	# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#-----------------------------------------------------------------------
# a couple of test stopwords to test that the words are really being
# configured from this file:
stopworda
stopwordb

#Standard english stop words taken from Lucene's StopAnalyzer
a
an
and
are
as
at
be
but
by
for
if
in
into
is
it
no
not
of
on
or
s
such
t
that
the
their
then
there
these
they
this
to
was
will
with
au
aux
avec
ce
ces
dans
de
des
du
elle
en
et
eux
il
je
la
le
leur
lui
ma
mais
me
même
mes
moi
mon
ne
nos
notre
nous
on
ou
par
pas
pour
qu
que
qui
sa
se
ses
son
sur
ta
te
tes
toi
ton
tu
un
une
vos
votre
vous
c
d
j
l
à
m
n
s
t
y
été
étée
étées
étés
étant
suis
es
est
sommes
êtes
sont
serai
seras
sera
serons
serez
seront
serais
serait
serions
seriez
seraient
étais
était
étions
étiez
étaient
fus
fut
fûmes
fûtes
furent
sois
soit
soyons
soyez
soient
fusse
fusses
fût
fussions
fussiez
fussent
ayant
eu
eue
eues
eus
ai
as
avons
avez
ont
aurai
auras
aura
aurons
aurez
auront
aurais
aurait
aurions
auriez
auraient
avais
avait
avions
aviez
avaient
eut
eûmes
eûtes
eurent
aie
aies
ait
ayons
ayez
aient
eusse
eusses
eût
eussions
eussiez
eussent
ceci
celà
cet
cette
ici
ils
les
leurs
quel
quels
quelle
quelles
sans
soi

`,
	"nl": `
aan
af
al
alleen
als
altijd
ben
bent
bij
daar
dag
dan
dat
de
der
deze
die
direct
dit
doch
doen
dus
een
eens
en
er
gaan
gaat
ge
geen
geweest
graag
haar
had
heb
hebben
heeft
hem
het
hij
hoe
hun
ik
in
is
je
kan
komt
kon
kunnen
kunt
laatste
maar
maken
me
mee
meer
men
met
mij
mijn
na
naar
niet
nog
nu
of
om
omdat
onder
ons
onze
ook
op
reeds
te
toch
toen
tot
uit
uw
van
vanaf
veel
via
voor
waar
was
wat
we
weer
wel
werd
wie
wij
wilt
worden
wordt
zal
ze
zei
zelf
zich
zij
zijn
zo
zoals
zou
`,
	"sv": `
#-----------------------------------------------------------------------
# translated
#-----------------------------------------------------------------------

kunna
om
ovan
enligt
i enlighet med detta
över
faktiskt
efter
efteråt
igen
mot
är inte
alla
tillåta
tillåter
nästan
ensam
längs
redan
också
även om
alltid
am
bland
bland
en
och
en annan
någon
någon
hur som helst
någon
något
ändå
ändå
var som helst
isär
visas
uppskatta
lämpligt
är
inte
runt
som
åt sidan
be
frågar
associerad
vid
tillgängliga
bort
väldigt
vara
blev
eftersom
bli
blir
blir
varit
innan
förhand
bakom
vara
tro
nedan
bredvid
förutom
bäst
bättre
mellan
bortom
både
kort
men
genom
c
c'mon
c: s
kom
kampanj
kan
kan inte
kan inte
cant
orsaka
orsaker
viss
säkerligen
förändringar
klart
co
com
komma
kommer
om
följaktligen
överväga
överväger
innehålla
innehållande
innehåller
motsvarande
kunde
kunde inte
kurs
närvarande
definitivt
beskrivits
trots
gjorde
inte
olika
göra
gör
inte
gör
inte
gjort
ned
nedåt
under
varje
edu
åtta
antingen
annars
någon annanstans
tillräckligt
godkändes
helt
speciellt
et
etc
även
någonsin
varje
alla
alla
allt
överallt
ex
exakt
exempel
utom
långt
få
femte
först
finansiella
fem
följt
efter
följer
för
fd
tidigare
framåt
fyra
från
ytterligare
dessutom
få
blir
få
given
ger
gå
går
gå
borta
fick
fått
hälsningar
hade
hade inte
händer
knappast
har
har inte
ha
har inte
med
han
han är
hallå
hjälpa
hence
henne
här
här finns
härefter
härmed
häri
härpå
hennes
själv
hej
honom
själv
hans
hit
förhoppningsvis
hur
howbeit
dock
jag skulle
jag ska
jag är
jag har
om
ignoreras
omedelbar
i
eftersom
inc
indeed
indikera
indikerade
indikerar
inre
mån
istället
in
inåt
är
är inte
den
det skulle
det ska
det är
dess
själv
bara
hålla
håller
hålls
vet
vet
känd
sista
nyligen
senare
senare
latterly
minst
mindre
lest
låt
låt oss
liknande
gillade
sannolikt
lite
ser
ser
ser
ltd
huvudsakligen
många
kan
kanske
mig
betyda
under tiden
endast
kanske
mer
dessutom
mest
mestadels
mycket
måste
min
själv
namn
nämligen
nd
nära
nästan
nödvändigt
behöver
behov
varken
aldrig
ändå
ny
nästa
nio
ingen
ingen
icke
ingen
ingen
eller
normalt
inte
ingenting
roman
nu
ingenstans
uppenbarligen
av
off
ofta
oh
ok
okay
gammal
på
en gång
ett
ettor
endast
på
eller
andra
andra
annars
borde
vår
vårt
oss
ut
utanför
över
övergripande
egen
särskilt
särskilt
per
kanske
placeras
vänligen
plus
möjligt
förmodligen
förmodligen
ger
ganska
citera
kvartalsvis
snarare
verkligen
rimligen
om
oavsett
gäller
relativt
respektive
höger
sa
samma
såg
säga
säger
säger
andra
det andra
se
ser
verkar
verkade
informationsproblem
verkar
sett
själv
själva
förnuftig
skickas
allvarlig
allvarligt
sju
flera
skall
hon
bör
bör inte
eftersom
sex
så
några
någon
på något sätt
någon
något
sometime
ibland
något
någonstans
snart
sorry
specificerade
ange
ange
fortfarande
sub
sådan
sup
säker
t s
ta
tas
berätta
tenderar
än
tacka
tack
thanx
att
det är
brinner
den
deras
deras
dem
själva
sedan
därifrån
där
det finns
därefter
därigenom
därför
däri
theres
därpå
dessa
de
de hade
de kommer
de är
de har
tror
tredje
detta
grundlig
grundligt
de
though
tre
genom
hela
thru
sålunda
till
tillsammans
alltför
tog
mot
mot
försökte
försöker
verkligt
försök
försöker
två gånger
två
enligt
tyvärr
såvida inte
osannolikt
tills
åt
upp
på
oss
använda
används
användbar
använder
användning
vanligtvis
uucp
värde
olika
mycket
via
viz
vs
vill
vill
var
var inte
sätt
vi
vi skulle
vi kommer
vi är
vi har
välkommen
väl
gick
var
var inte
vad
vad är
oavsett
när
varifrån
närhelst
där
var är
varefter
medan
varigenom
vari
varpå
varhelst
huruvida
som
medan
dit
som
vem är
vem
hela
vem
vars
varför
kommer
villig
önskar
med
inom
utan
kommer inte
undrar
skulle
skulle inte
ja
ännu
ni
du skulle
kommer du
du är
du har
din
själv
er
noll
tjänsteman
skarpt
kritiserade
`,
	"zh": `
的
一
不
在
人
有
是
为
以
于
上
他
而
后
之
来
及
了
因
下
可
到
由
这
与
也
此
但
并
个
其
已
无
小
我
们
起
最
再
今
去
好
只
又
或
很
亦
某
把
那
你
乃
它
吧
被
比
别
趁
当
从
到
得
打
凡
儿
尔
该
各
给
跟
和
何
还
即
几
既
看
据
距
靠
啦
了
另
么
每
们
嘛
拿
哪
那
您
凭
且
却
让
仍
啥
如
若
使
谁
虽
随
同
所
她
哇
嗡
往
哪
些
向
沿
哟
用
于
咱
则
怎
曾
至
致
着
诸
自
`,
	"ru": `
а
е
и
ж
м
о
на
не
ни
об
но
он
мне
мои
мож
она
они
оно
мной
много
многочисленное
многочисленная
многочисленные
многочисленный
мною
мой
мог
могут
можно
может
можхо
мор
моя
моё
мочь
над
нее
оба
нам
нем
нами
ними
мимо
немного
одной
одного
менее
однажды
однако
меня
нему
меньше
ней
наверху
него
ниже
мало
надо
один
одиннадцать
одиннадцатый
назад
наиболее
недавно
миллионов
недалеко
между
низко
меля
нельзя
нибудь
непрерывно
наконец
никогда
никуда
нас
наш
нет
нею
неё
них
мира
наша
наше
наши
ничего
начала
нередко
несколько
обычно
опять
около
мы
ну
нх
от
отовсюду
особенно
нужно
очень
отсюда
в
во
вон
вниз
внизу
вокруг
вот
восемнадцать
восемнадцатый
восемь
восьмой
вверх
вам
вами
важное
важная
важные
важный
вдали
везде
ведь
вас
ваш
ваша
ваше
ваши
впрочем
весь
вдруг
вы
все
второй
всем
всеми
времени
время
всему
всего
всегда
всех
всею
всю
вся
всё
всюду
г
год
говорил
говорит
года
году
где
да
ее
за
из
ли
же
им
до
по
ими
под
иногда
довольно
именно
долго
позже
более
должно
пожалуйста
значит
иметь
больше
пока
ему
имя
пор
пора
потом
потому
после
почему
почти
посреди
ей
два
две
двенадцать
двенадцатый
двадцать
двадцатый
двух
его
дел
или
без
день
занят
занята
занято
заняты
действительно
давно
девятнадцать
девятнадцатый
девять
девятый
даже
алло
жизнь
далеко
близко
здесь
дальше
для
лет
зато
даром
первый
перед
затем
зачем
лишь
десять
десятый
ею
её
их
бы
еще
при
был
про
процентов
против
просто
бывает
бывь
если
люди
была
были
было
будем
будет
будете
будешь
прекрасно
буду
будь
будто
будут
ещё
пятнадцать
пятнадцатый
друго
другое
другой
другие
другая
других
есть
пять
быть
лучше
пятый
к
ком
конечно
кому
кого
когда
которой
которого
которая
которые
который
которых
кем
каждое
каждая
каждые
каждый
кажется
как
какой
какая
кто
кроме
куда
кругом
с
т
у
я
та
те
уж
со
то
том
снова
тому
совсем
того
тогда
тоже
собой
тобой
собою
тобою
сначала
только
уметь
тот
тою
хорошо
хотеть
хочешь
хоть
хотя
свое
свои
твой
своей
своего
своих
свою
твоя
твоё
раз
уже
сам
там
тем
чем
сама
сами
теми
само
рано
самом
самому
самой
самого
семнадцать
семнадцатый
самим
самими
самих
саму
семь
чему
раньше
сейчас
чего
сегодня
себе
тебе
сеаой
человек
разве
теперь
себя
тебя
седьмой
спасибо
слишком
так
такое
такой
такие
также
такая
сих
тех
чаще
четвертый
через
часто
шестой
шестнадцать
шестнадцатый
шесть
четыре
четырнадцать
четырнадцатый
сколько
сказал
сказала
сказать
ту
ты
три
эта
эти
что
это
чтоб
этом
этому
этой
этого
чтобы
этот
стал
туда
этим
этими
рядом
тринадцать
тринадцатый
этих
третий
тут
эту
суть
чуть
тысяч`,
	"sr": `
# Serbian Cyrillic - Српски Ћирилица

а
ако
али
баш
без
би
биће
бих
била
били
било
био
бисмо
бисте
бити
близу
број
ће
ћемо
ћеш
често
ћете
чији
ћу
да
дана
данас
до
добар
добити
доћи
док
доле
дошао
други
дуж
два
га
где
горе
хоће
хоћемо
хоћеш
хоћете
хоћу
хвала
и
иако
ићи
иде
их
или
има
имам
имао
испод
из
између
изнад
изван
изволи
ја
је
један
једини
једна
једне
једно
једном
јер
јесам
јеси
јесмо
јесте
јесу
јој
још
јуче
кад
када
како
као
кога
која
које
који
којима
коју
кроз
ли
мали
мањи
ме
мене
мени
ми
мимо
мисли
много
моћи
могу
мој
моја
моје
мора
морао
на
наћи
над
након
нам
нама
нас
наш
наша
наше
нашег
не
неће
нећемо
нећеш
нећете
нећу
негде
него
нека
некад
неки
неко
неког
неку
нема
немам
нешто
ни
није
ниједан
никада
никога
нисам
ниси
нисмо
ништа
нисте
нису
њега
његов
његова
његово
њему
њен
њих
њихов
њихова
њихово
њим
њима
њој
њу
о
од
одмах
око
около
он
она
онај
они
оно
осим
остали
отишао
ова
овако
овамо
овде
ове
ово
па
питати
по
почетак
под
поједини
поред
после
поводом
правити
пре
преко
према
први
пут
радије
са
сада
сам
само
се
себе
себи
си
смети
смо
шта
сте
што
ствар
стварно
су
сутра
сваки
све
сви
свим
свог
свој
своја
своје
свом
свугде
та
тачно
тада
тај
тако
такође
тамо
те
тебе
теби
ти
тим
то
тој
томе
ту
твој
твоја
твоје
у
учинио
учинити
умало
унутра
употребити
уз
узети
вам
вама
вас
ваш
ваша
ваше
већ
већина
веома
ви
више
врло
за
захвалити
зар
зашто
због
желео
жели
знати

# Serbian Latin - Srpski Latinica

a
ako
ali
baš
bez
bi
biće
bih
bila
bili
bilo
bio
bismo
biste
biti
blizu
broj
će
ćemo
ćeš
često
ćete
čiji
ću
da
dana
danas
do
dobar
dobiti
doći
dok
dole
došao
drugi
duž
dva
ga
gde
gore
hoće
hoćemo
hoćeš
hoćete
hoću
hvala
i
iako
ići
ide
ih
ili
ima
imam
imao
ispod
iz
između
iznad
izvan
izvoli
ja
je
jedan
jedini
jedna
jedne
jedno
jednom
jer
jesam
jesi
jesmo
jest
jeste
jesu
joj
još
juče
kad
kada
kako
kao
koga
koja
koje
koji
kojima
koju
kroz
li
mali
manji
me
mene
meni
mi
mimo
misli
mnogo
moći
mogu
moj
moja
moje
mora
morao
na
naći
nad
nakon
nam
nama
nas
naš
naša
naše
našeg
ne
neće
nećemo
nećeš
nećete
neću
negde
nego
neka
nekad
neki
neko
nekog
neku
nema
nemam
nešto
ni
nije
nijedan
nikada
nikoga
nisam
nisi
nismo
ništa
niste
nisu
njega
njegov
njegova
njegovo
njemu
njen
njih
njihov
njihova
njihovo
njim
njima
njoj
nju
o
od
odmah
oko
okolo
on
ona
onaj
oni
ono
osim
ostali
otišao
ova
ovako
ovamo
ovde
ove
ovo
pa
pitati
po
početak
pod
pojedini
pored
posle
povodom
praviti
pre
preko
prema
prvi
put
radije
sa
sada
sam
samo
se
sebe
sebi
si
smeti
smo
šta
ste
što
stvar
stvarno
su
sutra
svaki
sve
svi
svim
svog
svoj
svoja
svoje
svom
svugde
ta
tačno
tada
taj
tako
takođe
tamo
te
tebe
tebi
ti
tim
to
toj
tome
tu
tvoj
tvoja
tvoje
u
učinio
učiniti
umalo
unutra
upotrebiti
uz
uzeti
vam
vama
vas
vaš
vaša
vaše
već
većina
veoma
vi
više
vrlo
za
zahvaliti
zar
zašto
zbog
želeo
želi
znati`,
	"bg": `
# This file was created by Jacques Savoy and is distributed under the BSD license.
# See http://members.unine.ch/jacques.savoy/clef/index.html.
# Also see http://www.opensource.org/licenses/bsd-license.html
а
аз
ако
ала
бе
без
беше
би
бил
била
били
било
близо
бъдат
бъде
бяха
в
вас
ваш
ваша
вероятно
вече
взема
ви
вие
винаги
все
всеки
всички
всичко
всяка
във
въпреки
върху
г
ги
главно
го
д
да
дали
до
докато
докога
дори
досега
доста
е
едва
един
ето
за
зад
заедно
заради
засега
затова
защо
защото
и
из
или
им
има
имат
иска
й
каза
как
каква
какво
както
какъв
като
кога
когато
което
които
кой
който
колко
която
къде
където
към
ли
м
ме
между
мен
ми
мнозина
мога
могат
може
моля
момента
му
н
на
над
назад
най
направи
напред
например
нас
не
него
нея
ни
ние
никой
нито
но
някои
някой
няма
обаче
около
освен
особено
от
отгоре
отново
още
пак
по
повече
повечето
под
поне
поради
после
почти
прави
пред
преди
през
при
пък
първо
с
са
само
се
сега
си
скоро
след
сме
според
сред
срещу
сте
съм
със
също
т
тази
така
такива
такъв
там
твой
те
тези
ти
тн
то
това
тогава
този
той
толкова
точно
трябва
тук
тъй
тя
тях
у
харесва
ч
че
често
чрез
ще
щом
я
`,
	"da": `
# next languages added from the python goose
# https://github.com/grangier/python-goose/tree/develop/goose/resources/text
af
alle
andet
andre
at
begge
da
de
den
denne
der
deres
det
dette
dig
din
dog
du
ej
eller
en
end
ene
eneste
enhver
et
fem
fire
flere
fleste
for
fordi
forrige
fra
få
før
god
han
hans
har
hendes
her
hun
hvad
hvem
hver
hvilken
hvis
hvor
hvordan
hvorfor
hvornår
i
ikke
ind
ingen
intet
jeg
jeres
kan
kom
kommer
lav
lidt
lille
man
mand
mange
med
meget
men
mens
mere
mig
ned
ni
nogen
noget
ny
nyt
nær
næste
næsten
og
op
otte
over
på
se
seks
ses
som
stor
store
syv
ti
til
to
tre
ud
var
`,
	"de": `
der
die
und
in
den
von
zu
mit
ist
das
des
im
für
auf
sich
dem
Die
nicht
ein
eine
als
auch
an
es
er
aus
bei
werden
sie
nach
Der
sind
war
wurde
wird
einer
Das
hat
am
wie
um
Sie
zum
oder
einen
über
dass
einem
noch
bis
nur
vor
zur
durch
so
haben
aber
ich
In
man
mehr
wir
daß
kann
sein
vom
Es
unter
Ich
hatte
gegen
Im
Er
wenn
dieser
seine
eines
können
diese
wieder
wurden
dann
was
schon
Jahr
zwei
seiner
Jahre
Jahren
ihre
gibt
zwischen
Ein
immer
waren
Zeit
Uhr
keine
Wir
sei
habe
sehr
hier
alle
Nach
ab
sowie
da
beim
heute
seit
diesem
uns
soll
Und
Deutschland
Mit
anderen
jedoch
ihr
damit
ersten
drei
Auch
doch
ihm
seinen
Stadt
etwa
sagte
ihn
Eine
sondern
bereits
müssen
ohne
Menschen
will
Prozent
ihrer
worden
Bei
selbst
jetzt
of
Als
seinem
neue
muss
allem
neuen
Ende
nun
Von
geht
ihren
SPD
So
Für
weil
wo
mich
mir
Aber
Am
Diese
ganz
dieses
etwas
andere
Geschichte
Frau
liegt
Wenn
ins
gut
einmal
konnte
Euro
du
denn
viele
Auf
machen
Herr
Leben
the
diesen
erst
lassen
Wie
dort
beiden
erste
The
Teil
deutschen
weiter
also
viel
sollte
dabei
Millionen
Was
später
hatten
während
Welt
ISBN
sagt
denen
wollen
steht
Da
kommt
kein
vier
nichts
de
allerdings
Seite
ob
dazu
gab
s
letzten
kam
USA
wegen
dies
zurück
großen
kommen
alles
rund
ja
sollen
deren
dafür
Doch
Kinder
wäre
Frage
weitere
würde
dessen
große
Januar
zwar
darauf
Arbeit

Beispiel
September
zusammen
einige
Land
allen
fast
Frauen
März
Namen
Unternehmen
ihrem
davon
Mann
Mai
Platz
deutsche
werde
Oktober
muß
Literatur
Art
ihnen
Deutschen
fünf
gilt
sehen
könnte
Dezember
stehen
sogar
seien
Wer
Seit
August
bin
Beifall
Fall
Juni
eigenen
November
mal
Film
finden
sagen
Regierung
April
München
oft
Dies
lange
ebenfalls
bekannt
Präsident
wohl
CDU/CSU
Zu
gehört
Man
weniger
gerade
statt
aller
Juli
möchte
Weg
Entwicklung
zunächst
ging
Mark
Bild
möglich
gar
besonders
hätte
macht
Politik
geben
Tag
Ihnen
Februar
Hier
Gemeinde
wenig
gewesen
Europa
gehen
gemacht
welche
New
gegenüber
heißt
Familie
Union
tun
Jahrhundert
einfach
Frankfurt
deutlich
Dabei
neben
sollten
Kirche
keinen
Artikel
Ihre
Peter
Thema
besteht
vielen
nie
bzw.
Aus
Zeitung
wollte
Kommission
seines
Hamburg
hätten
Geld
meine
Dr
kaum
zweiten
Während
lässt
Anfang
Um
Ort
weiß
findet
Bereich
Haus
anderem
Mal

deshalb
alten
erhalten
zehn
Zum
bisher
meisten
darüber
würden
hin
Form
An
bleibt
sieht
Gesellschaft
Berliner
Den
vergangenen
bezeichnet
Nr.
Ziel
je
weit
Grund
sechs
darf
Rolle
Deutsche
wissen
jeder
zeigt
Damit
Denn
mehrere
nächsten
Vor
Dann
schließlich
kleinen
Durch
Michael
km
Lage
Gruppe
Band
damals
Spiel
Sohn
Dr.
stark
Universität
Hilfe
besser
hinter
meist
Seine
St.
stellt
Tage
unsere
daher
Nur
wirklich
führt
Dieser
beispielsweise
kurz
Bericht
gleich
weiteren
Straße
bleiben
Wirtschaft
Siehe
Zukunft
eher
Bedeutung
Recht
insbesondere
Bevölkerung
schnell
nehmen
Verlag
CDU
Tod
Alle
solche
neu
Bundesregierung
pro
Frankreich
Jahres
konnten
Ihr
ließ
Du
kleine
Europäischen
Vater
genannt
lang
Titel
Rahmen
Wort
eigentlich
erhielt
einigen
Woche
FC
Musik
dagegen
Sein
allein
Einsatz
genau
begann
innerhalb
unserer
Partei
Polizei
Wasser
bringen
deutscher
natürlich
eigene

Wochen
insgesamt
Außerdem
Bis
halten
politischen
musste
Parlament
Meter
Hand
Zahl
stellen
gesagt
führen
daran
Erfolg
befindet
Zur
verschiedenen
Probleme
Unter
Abgeordneten
Milliarden
nahm
stand
geworden
c
liegen
erstmals
Sprache
Fragen
nämlich
Ja
Kollegen
Männer
Nicht
Wolfgang
Problem
Mutter
Minuten
Weitere
Mitte
Mitglied
Jahrhunderts
Krieg
Hans
könnten
Thomas
Über
Personen
Friedrich
ca.
ebenso
machte
York
vielleicht
Stelle
derzeit
Ländern
Höhe
verwendet
gute
überhaupt
Länder
Angaben
führte
gegeben
Tel.
klar
Karl
europäischen
sicher
Saison
Programm
erreicht
GRÜNEN
beide
Sonntag
sowohl
Region
alte
Staaten
Paris
Beginn
Buch
zweite
ganze
hinaus
König
Morgen
handelt
fand
Schweiz
jeweils
Weise
DM
fest
per
blieb
Mitglieder
Richtung
Heute
Stunden
leicht
Leute
wobei
gehören
bietet
Wien
politische
Folge
Blick
aufgrund
Entscheidung
Dort
Neben
hält
Gebiet
gemeinsam
erklärt
direkt
könne
Daten
recht
schwer
Bayern
jeden
Name
Schule
GmbH
dürfen
laut
Seiten
Bürger
Eltern
dpa
Meinung
Werke
Jetzt
letzte
Spieler
bald
London
häufig
heutigen
Einwohner
acht
eben
Internet
Markt
dich
Nein
Situation
System
zuvor
Möglichkeit
Freitag
mein
Mannheim
Fenster
Kosten
inzwischen
kamen
John
sieben
bekommen
erreichen
unser
Verfügung
Köln
Dazu
besten
Zusammenhang
Reihe
Kritik
richtig
Liste
Herren
Augen
taz
zeigen
siehe
hohen
spielte
leben
völlig
Neue
ihres
spielt
Sicherheit
weiterhin
hoch
nachdem
gegründet
erneut
sah
z.
wer
Informationen
anders
spielen
Dieses
gleichen
Kultur
größten
eingesetzt
Unterstützung
Beim
erklärte
Allerdings
Firma
Amt
Kopf
trotz
Erst
gebracht
gestellt
läuft
schließen
Bilder
nimmt
Mitarbeiter
BÜNDNIS
Deshalb
verschiedene
zudem
Werk
Ergebnis
Heinrich
Bau
ehemaligen
Preis
Tochter
Stuttgart
Samstag
Bad
Verfahren
Kind
früher
Paul
darin
paar
Punkt
Weblinks
Nun
Maßnahmen
Österreich
Wilhelm
Herrn
z.B.
Noch
Staat
Zusammenarbeit
knapp
Nacht
einzelnen
trat
gestern
Team
Osten
scheint
Mannschaft
Tagen
internationalen
jede
mindestens
teilweise
einzige
Soldaten
setzt
gefunden
Kunst
lediglich
öffentlichen
bedeutet
Raum
gewann
Kampf
Martin
Ist
Begriff
Hause
entwickelt
Wahl
Schon
arbeiten
größte
Donnerstag
Ab
Viele
Quellen
Nachdem
dadurch
Italien
erster
gekommen
dir
Mittwoch
danach
stellte
her
zahlreiche
Landes
Gesetz
Monaten
PDS
Rat
Franz
Verein
sonst
Frankfurter
Meine
Klaus
Karriere
müsse
meiner
anderer
zuletzt
Monate
Alter
hohe
Interesse
Regie
Montag
genommen
lag
Sommer
spricht
Trainer
Liebe
jedem
/DIE
Westen
guten
Kilometer
Johann
gesehen
darunter
solchen
indem
Mittel
oben
Schweizer
wichtig
Hälfte
Regel
obwohl
Bürgermeister
Aufgabe
Spiele
folgenden
Dienstag
version
Sache
sprechen
Gemeinden
electronic
for
Norden
außerdem
Antrag
gleichzeitig
ganzen
Politiker
gehörte
großer
China
Nähe
bereit
setzte
Druck
tatsächlich
Gott
frei
Grünen
zumindest
Opfer
genug
versucht
bevor
`,
	"fi": `
alla
ansiosta
ehkä
ei
enemmän
ennen
etessa
f
haikki
he
hitaasti
hoikein
hyvin
hän
ilman
ja
jos
jälkeen
kanssa
kaukana
kenties
keskellä
kesken
koskaan
kuinkan
kukka
kylliksi
kyllä
liian
lla
lla
luona
lähellä
läpi
me
miksi
mikä
milloin
milloinkan
minä
missä
miten
nopeasti
nyt
oikea
oikealla
paljon
siellä
sinä
ssa
sta
suoraan
tai
takana
takia
tarpeeksi
te
tässä
ulkopuolella
vahemmän
vasen
vasenmalla
vastan
vielä
vieressä
vähän
yhdessä
ylös
`,
	"hu": `
a
á
ahogy
ahol
aki
akik
akkor
alatt
által
általában
amely
amelyek
amelyekben
amelyeket
amelyet
amelynek
ami
amit
amolyan
amp
amíg
amikor
át
abban
ahhoz
annak
arra
arról
az
azok
azon
azt
azzal
azért
aztán
azután
azonban
b
bár
be
belül
benne
c
cikk
cikkek
cikkeket
csak
d
de
e
é
eddig
egész
egy
egyes
egyetlen
egyéb
egyik
egyre
ekkor
el
elég
ellen
elő
először
előtt
első
én
éppen
ebben
ehhez
emilyen
ennek
erre
ez
ezt
ezek
ezen
ezzel
ezért
és
f
fel
felé
g
h
hanem
hiszen
hogy
hogyan
i
í
igen
így
illetve
ill.
ill
ilyen
ilyenkor
is
ison
ismét
itt
j
jó
jól
jobban
k
kell
kellett
keresztül
keressünk
ki
kívül
között
közül
l
legalább
lehet
lehetett
legyen
lenne
lenni
lesz
lett
m
maga
magát
majd
majd
már
más
másik
meg
még
mellett
mert
mely
melyek
mi
mit
míg
miért
milyen
mikor
minden
mindent
mindenki
mindig
mint
mintha
mivel
most
n
nagy
nagyobb
nagyon
ne
néha
nekem
neki
nem
néhány
nélkül
nincs
o
ó
olyan
ott
össze
ö
ő
ők
őket
p
pedig
persze
q
r
rá
s
saját
sem
semmi
sok
sokat
sokkal
sz
számára
szemben
szerint
szinte
t
talán
tehát
teljes
tovább
továbbá
több
u
ú
úgy
ugyanis
új
újabb
újra
után
utána
utolsó
ü
ű
v
vagy
vagyis
valaki
valamely
valami
valamint
való
vagyok
van
vannak
volt
voltam
voltak
voltunk
vissza
vele
viszont
volna
számolnak
szólnak
szól
w
x
y
z
zs
a
ahogy
ahol
aki
akkor
alatt
általában
által
amely
amíg
amikor
ami
amolyan
arra
át
az
azért
azonban
azon
aztán
azt
azután
azzal
bár
be
belül
benne
cikk
csak
de
eddig
egész
egy
egyéb
egyes
egyetlen
egyik
egyre
ekkor
el
elég
ellen
elő
először
előtt
első
emilyen
én
éppen
erre
és
e
ez
ezen
ezért
ezzel
fel
felé
hanem
hiszen
hogy
hogyan
igen
így
ill.
illetve
ill
ilyen
ilyenkor
ismét
ison
itt
jó
jobban
jól
kell
keres
keresztül
ki
kívül
között
közül
legalább
legyen
lehet
lenni
lett
maga
maga
majd
már
más
másik
még
meg
mellett
mely
mert
miért
míg
mikor
milyen
minden
mindenki
mindig
mi
mint
mintha
mivel
most
nagy
nagyobb
nagyon
ne
néha
néhány
neki
nélkül
nem
nincs
ők
olyan
ő
össze
ott
pedig
persze
rá
saját
s
sem
semmi
sokkal
sok
számára
számol
szemben
szerint
szinte
szól
talán
tehát
teljes
továbbá
tovább
úgy
ugyanis
új
újabb
újra
utána
után
utolsó
vagy
vagyis
valaki
valamely
valami
valamint
való
van
vissza
viszont
volt
`,
	"id": `
a
abad
acara
aceh
ada
adalah
adanya
adapun
agak
agaknya
agama
agar
agustus
air
akan
akankah
akhir
akhiri
akhirnya
akibat
aku
akulah
alam
album
amat
amatlah
amerika
anak
and
anda
andalah
anggota
antar
antara
antarabangsa
antaranya
apa
apaan
apabila
apakah
apalagi
apatah
api
april
artikel
artinya
as
asal
asalkan
asas
asia
asing
atas
atau
ataukah
ataupun
australia
awal
awalnya
awam
b
badan
bagai
bagaikan
bagaimana
bagaimanakah
bagaimanapun
bagainamakah
bagi
bagian
bahagian
bahan
baharu
bahasa
bahawa
bahkan
bahwa
bahwasannya
bahwasanya
baik
baiknya
bakal
bakalan
balik
bandar
bangsa
bank
banyak
bapak
barang
barangan
barat
baru
baru-baru
bawah
beberapa
begini
beginian
beginikah
beginilah
begitu
begitukah
begitulah
begitupun
bekas
bekerja
belakang
belakangan
belanda
beli
beliau
belum
belumlah
benar
benarkah
benarlah
bentuk
berada
berakhir
berakhirlah
berakhirnya
berapa
berapakah
berapalah
berapapun
berarti
berasal
berat
berawal
berbagai
berbanding
berbeda
berdasarkan
berdatangan
berharap
berhasil
beri
berikan
berikut
berikutan
berikutnya
berita
berjalan
berjaya
berjumlah
berkaitan
berkali
berkali-kali
berkata
berkehendak
berkeinginan
berkenaan
berlainan
berlaku
berlalu
berlangsung
berlebihan
bermacam
bermacam-macam
bermain
bermaksud
bermula
bernama
bernilai
bersama
bersama-sama
bersiap
bertanya
bertemu
berturut
bertutur
berubah
berujar
berupa
besar
besok
betul
betulkah
bhd
biasa
biasanya
bidang
bila
bilakah
bilion
bintang
bisa
bisakah
blog
bn
bola
boleh
bolehkah
bolehlah
buat
bukan
bukankah
bukanlah
bukannya
buku
bulan
bumi
bung
bursa
cadangan
cara
caranya
catch
china
click
code
copyright
cukup
cukupkah
cukuplah
cuma
daerah
dagangan
dahulu
dalam
dan
dana
dapat
dari
daripada
dasar
data
datang
datuk
dekat
demi
demikian
demikianlah
dengan
depan
derivatives
desa
desember
detik
dewan
di
dia
diadakan
diakhiri
diakhirinya
dialah
dianggap
diantara
diantaranya
diberi
diberikan
diberikannya
dibuat
dibuatnya
dibuka
dicatatkan
didapat
didatangkan
didirikan
diduga
digunakan
diibaratkan
diibaratkannya
diingat
diingatkan
diinginkan
dijangka
dijawab
dijelaskan
dijelaskannya
dikarenakan
dikatakan
dikatakannya
dikenal
dikerjakan
diketahui
diketahuinya
dikira
dilakukan
dilalui
dilihat
dimaksud
dimaksudkan
dimaksudkannya
dimaksudnya
dimana
diminta
dimintai
dimisalkan
dimulai
dimulailah
dimulainya
dimungkinkan
dini
diniagakan
dipastikan
diperbuat
diperbuatnya
dipergunakan
diperkirakan
diperlihatkan
diperlukan
diperlukannya
dipersoalkan
dipertanyakan
dipunyai
diri
dirilis
dirinya
dis
disampaikan
disebut
disebutkan
disebutkannya
disember
disini
disinilah
distrik
ditambahkan
ditandaskan
ditanya
ditanyai
ditanyakan
ditegaskan
ditemukan
ditujukan
ditunjuk
ditunjuki
ditunjukkan
ditunjukkannya
ditunjuknya
ditutup
dituturkan
dituturkannya
diucapkan
diucapkannya
diungkapkan
document.write
dolar
dong
dr
dua
dulu
dunia
effective
ekonomi
eksekutif
eksport
empat
enam
enggak
enggaknya
entah
entahlah
era
eropa
err
faedah
feb
film
gat
gedung
gelar
gettracker
global
grup
guna
gunakan
gunung
hadap
hadapan
hal
hampir
hanya
hanyalah
harga
hari
harian
harus
haruslah
harusnya
hasil
hendak
hendaklah
hendaknya
hidup
hingga
https
hubungan
hukum
hutan
i
ia
iaitu
ialah
ibarat
ibaratkan
ibaratnya
ibu
ii
iklan
ikut
ilmu
indeks
india
indonesia
industri
informasi
ingat
inggris
ingin
inginkah
inginkan
ini
inikah
inilah
internasional
islam
isnin
isu
italia
itu
itukah
itulah
jabatan
jadi
jadilah
jadinya
jakarta
jalan
jalur
jaman
jan
jangan
jangankan
janganlah
januari
jauh
jawa
jawab
jawaban
jawabnya
jawatan
jawatankuasa
jelas
jelaskan
jelaslah
jelasnya
jenis
jepang
jepun
jerman
jika
jikalau
jiwa
jual
jualan
juga
julai
jumaat
jumat
jumlah
jumlahnya
jun
juni
justru
juta
kabar
kabupaten
kadar
kala
kalangan
kalau
kalaulah
kalaupun
kali
kalian
kalimantan
kami
kamilah
kamis
kamu
kamulah
kan
kantor
kapal
kapan
kapankah
kapanpun
karena
karenanya
karya
kasus
kata
katakan
katakanlah
katanya
kaunter
kawasan
ke
keadaan
kebetulan
kebutuhan
kecamatan
kecil
kedua
kedua-dua
keduanya
kedudukan
kegiatan
kehidupan
keinginan
kejadian
kekal
kelamaan
kelihatan
kelihatannya
kelima
kelompok
keluar
keluarga
kelurahan
kembali
kementerian
kemudahan
kemudian
kemungkinan
kemungkinannya
kenaikan
kenapa
kenyataan
kepada
kepadanya
kepala
kepentingan
keputusan
kerajaan
kerana
kereta
kerja
kerjasama
kes
kesampaian
keselamatan
keseluruhan
keseluruhannya
kesempatan
kesihatan
keterangan
keterlaluan
ketiga
ketika
ketua
keuntungan
kewangan
khamis
khusus
khususnya
kini
kinilah
kira
kira-kira
kiranya
kita
kitalah
klci
klibor
klik
km
kok
komentar
kompas
komposit
kondisi
kontrak
korban
korea
kos
kota
kuala
kuasa
kukuh
kumpulan
kurang
kurangnya
lagi
lagian
lagu
lah
lain
lainnya
laku
lalu
lama
lamanya
langkah
langsung
lanjut
lanjutnya
laporan
laut
lebih
lembaga
lepas
lewat
lima
lingkungan
login
lokasi
lot
luar
luas
lumpur
mac
macam
mahkamah
mahu
majlis
maka
makanan
makanya
makin
maklumat
malah
malahan
malam
malaysia
mampu
mampukah
mana
manakala
manalagi
mantan
manusia
masa
masalah
masalahnya
masih
masihkah
masing
masing-masing
masuk
masyarakat
mata
mau
maupun
measure
media
mei
melainkan
melakukan
melalui
melawan
melihat
melihatnya
memandangkan
memang
memastikan
membantu
membawa
memberi
memberikan
membolehkan
membuat
memerlukan
memihak
memiliki
meminta
memintakan
memisalkan
memperbuat
mempergunakan
memperkirakan
memperlihatkan
mempersiapkan
mempersoalkan
mempertanyakan
mempunyai
memulai
memungkinkan
menaiki
menambah
menambahkan
menandaskan
menanti
menantikan
menanya
menanyai
menanyakan
menarik
menawarkan
mencapai
mencari
mencatatkan
mendapat
mendapatkan
mendatang
mendatangi
mendatangkan
menegaskan
menerima
menerusi
mengadakan
mengakhiri
mengaku
mengalami
mengambil
mengapa
mengatakan
mengatakannya
mengenai
mengerjakan
mengetahui
menggalakkan
menggunakan
menghadapi
menghendaki
mengibaratkan
mengibaratkannya
mengikut
mengingat
mengingatkan
menginginkan
mengira
mengucapkan
mengucapkannya
mengumumkan
mengungkapkan
mengurangkan
meninggal
meningkat
meningkatkan
menjadi
menjalani
menjawab
menjelang
menjelaskan
menokok
menteri
menuju
menunjuk
menunjuki
menunjukkan
menunjuknya
menurut
menuturkan
menyaksikan
menyampaikan
menyangkut
menyatakan
menyebabkan
menyebutkan
menyediakan
menyeluruh
menyiapkan
merasa
mereka
merekalah
merosot
merupakan
meski
meskipun
mesyuarat
metrotv
meyakini
meyakinkan
milik
militer
minat
minggu
minta
minyak
mirip
misal
misalkan
misalnya
mobil
modal
mohd
mudah
mula
mulai
mulailah
mulanya
muncul
mungkin
mungkinkah
musik
musim
nah
naik
nama
namun
nanti
nantinya
nasional
negara
negara-negara
negeri
new
niaga
nilai
nomor
noun
nov
november
numeral
numeralia
nya
nyaris
nyatanya
of
ogos
okt
oktober
olah
oleh
olehnya
operasi
orang
organisasi
pada
padahal
padanya
pagetracker
pagi
pak
paling
pameran
panjang
pantas
papan
para
paras
parlimen
partai
parti
particle
pasar
pasaran
password
pasti
pastilah
pasukan
paticle
pegawai
pejabat
pekan
pekerja
pelabur
pelaburan
pelancongan
pelanggan
pelbagai
peluang
pemain
pembangunan
pemberita
pembinaan
pemerintah
pemerintahan
pemimpin
pendapatan
pendidikan
penduduk
penerbangan
pengarah
pengeluaran
pengerusi
pengguna
penggunaan
pengurusan
peniaga
peningkatan
penting
pentingnya
per
perancis
perang
peratus
percuma
perdagangan
perdana
peringkat
perjanjian
perkara
perkhidmatan
perladangan
perlu
perlukah
perlunya
permintaan
pernah
perniagaan
persekutuan
persen
persidangan
persoalan
pertama
pertandingan
pertanyaan
pertanyakan
pertubuhan
pertumbuhan
perubahan
perusahaan
pesawat
peserta
petang
pihak
pihaknya
pilihan
pinjaman
polis
polisi
politik
pos
posisi
presiden
prestasi
produk
program
projek
pronomia
pronoun
proses
proton
provinsi
pt
pubdate
pukul
pula
pulau
pun
punya
pusat
rabu
radio
raja
rakan
rakyat
ramai
rantau
rasa
rasanya
rata
raya
rendah
republik
resmi
ribu
ringgit
root
ruang
rumah
rupa
rupanya
saat
saatnya
sabah
sabtu
sahaja
saham
saja
sajalah
sakit
salah
saling
sama
sama-sama
sambil
sampai
sampaikan
sana
sangat
sangatlah
sarawak
satu
sawit
saya
sayalah
sdn
se
sebab
sebabnya
sebagai
sebagaimana
sebagainya
sebagian
sebahagian
sebaik
sebaiknya
sebaliknya
sebanyak
sebarang
sebegini
sebegitu
sebelah
sebelum
sebelumnya
sebenarnya
seberapa
sebesar
sebetulnya
sebisanya
sebuah
sebut
sebutlah
sebutnya
secara
secukupnya
sedang
sedangkan
sedemikian
sedikit
sedikitnya
seenaknya
segala
segalanya
segera
segi
seharusnya
sehingga
seingat
sejak
sejarah
sejauh
sejenak
sejumlah
sekadar
sekadarnya
sekali
sekali-kali
sekalian
sekaligus
sekalipun
sekarang
sekaranglah
sekecil
seketika
sekiranya
sekitar
sekitarnya
sekolah
sektor
sekurang
sekurangnya
sekuriti
sela
selagi
selain
selaku
selalu
selama
selama-lamanya
selamanya
selanjutnya
selasa
selatan
selepas
seluruh
seluruhnya
semacam
semakin
semalam
semampu
semampunya
semasa
semasih
semata
semaunya
sementara
semisal
semisalnya
sempat
semua
semuanya
semula
sen
sendiri
sendirian
sendirinya
senin
seolah
seolah-olah
seorang
sepak
sepanjang
sepantasnya
sepantasnyalah
seperlunya
seperti
sepertinya
sepihak
sept
september
serangan
serantau
seri
serikat
sering
seringnya
serta
serupa
sesaat
sesama
sesampai
sesegera
sesekali
seseorang
sesi
sesuai
sesuatu
sesuatunya
sesudah
sesudahnya
setelah
setempat
setengah
seterusnya
setiap
setiausaha
setiba
setibanya
setidak
setidaknya
setinggi
seusai
sewaktu
siap
siapa
siapakah
siapapun
siaran
sidang
singapura
sini
sinilah
sistem
soal
soalnya
sokongan
sri
stasiun
suara
suatu
sudah
sudahkah
sudahlah
sukan
suku
sumber
sungai
supaya
surat
susut
syarikat
syed
tadi
tadinya
tahap
tahu
tahun
tak
tama
tambah
tambahnya
tampak
tampaknya
tampil
tan
tanah
tandas
tandasnya
tanggal
tanpa
tanya
tanyakan
tanyanya
tapi
tawaran
tegas
tegasnya
teknologi
telah
televisi
teman
tempat
tempatan
tempo
tempoh
tenaga
tengah
tentang
tentara
tentu
tentulah
tentunya
tepat
terakhir
terasa
terbaik
terbang
terbanyak
terbesar
terbuka
terdahulu
terdapat
terdiri
terhadap
terhadapnya
teringat
terjadi
terjadilah
terjadinya
terkait
terkenal
terkira
terlalu
terlebih
terletak
terlihat
termasuk
ternyata
tersampaikan
tersebut
tersebutlah
tertentu
tertuju
terus
terutama
testimoni
testimony
tetap
tetapi
the
tiada
tiap
tiba
tidak
tidakkah
tidaklah
tidaknya
tiga
tim
timbalan
timur
tindakan
tinggal
tinggi
tingkat
toh
tokoh
try
tun
tunai
tunjuk
turun
turut
tutur
tuturnya
tv
uang
ucap
ucapnya
udara
ujar
ujarnya
umum
umumnya
unescape
ungkap
ungkapnya
unit
universitas
untuk
untung
upaya
urus
usah
usaha
usai
user
utama
utara
var
versi
waduh
wah
wahai
wakil
waktu
waktunya
walau
walaupun
wang
wanita
warga
warta
wib
wilayah
wong
word
ya
yaitu
yakin
yakni
yang
zaman
`,
	"it": `
ad            
al            
allo          
ai            
agli          
all           
agl           
alla          
alle          
con           
col           
coi           
da            
dal           
dallo         
dai           
dagli         
dall          
dagl          
dalla         
dalle         
di            
del           
dello         
dei           
degli         
dell          
degl          
della         
delle         
in            
nel           
nello         
nei           
negli         
nell          
negl          
nella         
nelle         
su            
sul           
sullo         
sui           
sugli         
sull          
sugl          
sulla         
sulle         
per           
tra           
contro        
io            
tu            
lui           
lei           
noi           
voi           
loro          
mio           
mia           
miei          
mie           
tuo           
tua           
tuoi          
tue           
suo           
sua           
suoi          
sue           
nostro        
nostra        
nostri        
nostre        
vostro        
vostra        
vostri        
vostre        
mi            
ti            
ci            
vi            
lo            
la            
li            
le            
gli           
ne            
il            
un            
uno           
una           
ma            
ed            
se            
perchè        
perché
perche
anche         
come          
dov           
dove          
che           
chi           
cui           
non           
più           
piu
quale         
quanto        
quanti        
quanta        
quante        
quello        
quelli        
quella        
quelle        
questo        
questi        
questa        
queste        
si            
tutto         
tutti         
a             
c             
e             
i             
l             
o             
ho
hai
ha
abbiamo
avete
hanno
abbia
abbiate
abbiano
avrò
avro
avrai
avrà
avra
avremo
avrete
avranno
avrei
avresti
avrebbe
avremmo
avreste
avrebbero
avevo
avevi
aveva
avevamo
avevate
avevano
ebbi
avesti
ebbe
avemmo
aveste
ebbero
avessi
avesse
avessimo
avessero
avendo
avuto
avuta
avuti
avute
sono
sei
è
é
e
siamo
siete
sia
siate
siano
sarà
sarai
sarò
saro
saremo
sarete
saranno
sarei
saresti
sarebbe
saremmo
sareste
sarebbero
ero
eri
era
eravamo
eravate
erano
fui
fosti
fu
fummo
foste
furono
fossi
fosse
fossimo
fossero
essendo
faccio
fai
facciamo
fanno
faccia
facciate
facciano
farà
farai
farò
faremo
farete
faranno
farei
faresti
farebbe
faremmo
fareste
farebbero
facevo
facevi
faceva
facevamo
facevate
facevano
feci
facesti
fece
facemmo
faceste
fecero
facessi
facesse
facessimo
facessero
facendo
sto
stai
sta
stiamo
stanno
stia
stiate
stiano
starà
starai
starò
staremo
starete
staranno
starei
staresti
starebbe
staremmo
stareste
starebbero
stavo
stavi
stava
stavamo
stavate
stavano
stetti
stesti
stette
stemmo
steste
stettero
stessi
stesse
stessimo
stessero
stando
`,
	"ko": `
을
의
에
이
를
으로
은
는
가
로
하고
과
에서
도
와
이다
고
부터
까지
께
에는
이라고
만
라고
보다
에도
다
토록
에게
나
대로
에서는
이나
이며
요
든
으로써
같이
로는
밖에
과의
며
로부터
처럼
아
라
여
으로는
이고
에서의
이라는
만에
으로부터
에서도
와의
엔
만을
부터는
만의
야
까지의
과는
치고
과를
으로의
까지는
보다는
만이
에만
로의
`,
	"nb": `
alle
andre
arbeid
av
begge
bort
bra
bruke
da
denne
der
deres
det
din
disse
du
eller
en
ene
eneste
enhver
enn
er
et
folk
for
fordi
forsÛke
fra
fÅ
fÛr
fÛrst
gjorde
gjÛre
god
gÅ
ha
hadde
han
hans
hennes
her
hva
hvem
hver
hvilken
hvis
hvor
hvordan
hvorfor
ikke
inn
innen
kan
kunne
lage
lang
lik
like
makt
mange
med
meg
meget
men
mens
mer
mest
min
mye
mÅ
mÅte
navn
nei
ny
nÅ
nÅr
og
ogsÅ
om
opp
oss
over
part
punkt
pÅ
rett
riktig
samme
sant
si
siden
sist
skulle
slik
slutt
som
start
stille
tid
til
tilbake
tilstand
under
ut
uten
var
ved
verdi
vi
vil
ville
vite
vÅr
vÖre
vÖrt
Å
`,
	"no": `
at
av
de
den
der
det
du
en
er
et
for
fra
før
med
og
om
over
på
som
til
ved
år
alle
bare
ble
bort
bra
da
deg
dem
denne
dere
deres
det
dette
din
disse
dit
ditt
eller
ene
enn
er
et
ett
etter
for
fram
først
få
god
gå
ha
han
hans
har
her
hit
hun
hva
hvem
hver
ikke
inn
ja
jeg
kan
kom
kun
kunne
lage
lang
lik
like
man
mer
min
mot
mye
må
måte
ned
nei
noe
noen
ny
nå
når
også
opp
oss
seg
selv
si
siden
sin
sine
sist
skal
skulle
slik
som
så
sånn
tid
til
under
ut
uten
var
ved
vi
vil
vite
vår
å
dei
di
då
eg
`,
	"pl": `
a
aby
ach
acz
aczkolwiek
aj
albo
ale
ależ
ani
aż
bardziej
bardzo
bo
bowiem
by
byli
bynajmniej
być
był
była
było
były
będzie
będą
cali
cała
cały
ci
cię
ciebie
co
cokolwiek
coś
czasami
czasem
czemu
czy
czyli
daleko
dla
dlaczego
dlatego
do
dobrze
dokąd
dość
dużo
dwa
dwaj
dwie
dwoje
dziś
dzisiaj
gdy
gdyby
gdyż
gdzie
gdziekolwiek
gdzieś
i
ich
ile
im
inna
inne
inny
innych
iż
ja
ją
jak
jakaś
jakby
jaki
jakichś
jakie
jakiś
jakiż
jakkolwiek
jako
jakoś
je
jeden
jedna
jedno
jednak
jednakże
jego
jej
jemu
jest
jestem
jeszcze
jeśli
jeżeli
już
ją
każdy
kiedy
kilka
kimś
kto
ktokolwiek
ktoś
która
które
którego
której
który
których
którym
którzy
ku
lat
lecz
lub
ma
mają
mało
mam
mi
mimo
między
mną
mnie
mogą
moi
moim
moja
moje
może
możliwe
można
mój
mu
musi
my
na
nad
nam
nami
nas
nasi
nasz
nasza
nasze
naszego
naszych
natomiast
natychmiast
nawet
nią
nic
nich
nie
niech
niego
niej
niemu
nigdy
nim
nimi
niż
no
o
obok
od
około
on
ona
one
oni
ono
oraz
oto
owszem
pan
pana
pani
po
pod
podczas
pomimo
ponad
ponieważ
powinien
powinna
powinni
powinno
poza
prawie
przecież
przed
przede
przedtem
przez
przy
roku
również
sam
sama
są
się
skąd
sobie
sobą
sposób
swoje
ta
tak
taka
taki
takie
także
tam
te
tego
tej
temu
ten
teraz
też
to
tobą
tobie
toteż
trzeba
tu
tutaj
twoi
twoim
twoja
twoje
twym
twój
ty
tych
tylko
tym
u
w
wam
wami
was
wasz
wasza
wasze
we
według
wiele
wielu
więc
więcej
wszyscy
wszystkich
wszystkie
wszystkim
wszystko
wtedy
wy
właśnie
z
za
zapewne
zawsze
ze
zł
znowu
znów
został
żaden
żadna
żadne
żadnych
że
żeby
`,
	"pt": `
último
é
acerca
agora
algmas
alguns
ali
ambos
antes
apontar
aquela
aquelas
aquele
aqueles
aqui
atrás
bem
bom
cada
caminho
cima
com
como
comprido
conhecido
corrente
das
debaixo
dentro
desde
desligado
deve
devem
deverá
direita
diz
dizer
dois
dos
e
ela
ele
eles
em
enquanto
então
está
estão
estado
estar
estará
este
estes
esteve
estive
estivemos
estiveram
eu
fará
faz
fazer
fazia
fez
fim
foi
fora
horas
iniciar
inicio
ir
irá
ista
iste
isto
ligado
maioria
maiorias
mais
mas
mesmo
meu
muito
muitos
nós
não
nome
nosso
novo
o
onde
os
ou
outro
para
parte
pegar
pelo
pessoas
pode
poderá
podia
por
porque
povo
promeiro
quê
qual
qualquer
quando
quem
quieto
são
saber
sem
ser
seu
somente
têm
tal
também
tem
tempo
tenho
tentar
tentaram
tente
tentei
teu
teve
tipo
tive
todos
trabalhar
trabalho
tu
um
uma
umas
uns
usa
usar
valor
veja
ver
verdade
verdadeiro
você
`,
	"hy": `
# next languages added from https://www.ranks.nl/stopwords
այդ
այլ
այն
այս
դու
դուք
եմ
են
ենք
ես
եք
է
էի
էին
էինք
էիր
էիք
էր
ըստ
թ
ի
ին
իսկ
իր
կամ
համար
հետ
հետո
մենք
մեջ
մի
ն
նա
նաև
նրա
նրանք
որ
որը
որոնք
որպես
ու
ում
պիտի
վրա
և
`,
	"eu": `
al
anitz
arabera
asko
baina
bat
batean
batek
bati
batzuei
batzuek
batzuetan
batzuk
bera
beraiek
berau
berauek
bere
berori
beroriek
beste
bezala
da
dago
dira
ditu
du
dute
edo
egin
ere
eta
eurak
ez
gainera
gu
gutxi
guzti
haiei
haiek
haietan
hainbeste
hala
han
handik
hango
hara
hari
hark
hartan
hau
hauei
hauek
hauetan
hemen
hemendik
hemengo
hi
hona
honek
honela
honetan
honi
hor
hori
horiei
horiek
horietan
horko
horra
horrek
horrela
horretan
horri
hortik
hura
izan
ni
noiz
nola
non
nondik
nongo
nor
nora
ze
zein
zen
zenbait
zenbat
zer
zergatik
ziren
zituen
zu
zuek
zuen
zuten
`,
	"bn": `
অবশ্য
অনেক
অনেকে
অনেকেই
অন্তত
অথবা
অথচ
অর্থাত
অন্য
আজ
আছে
আপনার
আপনি
আবার
আমরা
আমাকে
আমাদের
আমার
আমি
আরও
আর
আগে
আগেই
আই
অতএব
আগামী
অবধি
অনুযায়ী
আদ্যভাগে
এই
একই
একে
একটি
এখন
এখনও
এখানে
এখানেই
এটি
এটা
এটাই
এতটাই
এবং
একবার
এবার
এদের
এঁদের
এমন
এমনকী
এল
এর
এরা
এঁরা
এস
এত
এতে
এসে
একে
এ
ঐ
 ই
ইহা
ইত্যাদি
উনি
উপর
উপরে
উচিত
ও
ওই
ওর
ওরা
ওঁর
ওঁরা
ওকে
ওদের
ওঁদের
ওখানে
কত
কবে
করতে
কয়েক
কয়েকটি
করবে
করলেন
করার
কারও
করা
করি
করিয়ে
করার
করাই
করলে
করলেন
করিতে
করিয়া
করেছিলেন
করছে
করছেন
করেছেন
করেছে
করেন
করবেন
করায়
করে
করেই
কাছ
কাছে
কাজে
কারণ
কিছু
কিছুই
কিন্তু
কিংবা
কি
কী
কেউ
কেউই
কাউকে
কেন
কে
কোনও
কোনো
কোন
কখনও
ক্ষেত্রে
খুব
গুলি
গিয়ে
গিয়েছে
গেছে
গেল
গেলে
গোটা
চলে
ছাড়া
ছাড়াও
ছিলেন
ছিল
জন্য
জানা
ঠিক
তিনি
তিনঐ
তিনিও
তখন
তবে
তবু
তাঁদের
তাঁাহারা
তাঁরা
তাঁর
তাঁকে
তাই
তেমন
তাকে
তাহা
তাহাতে
তাহার
তাদের
তারপর
তারা
তারৈ
তার
তাহলে
তিনি
তা
তাও
তাতে
তো
তত
তুমি
তোমার
তথা
থাকে
থাকা
থাকায়
থেকে
থেকেও
থাকবে
থাকেন
থাকবেন
থেকেই
দিকে
দিতে
দিয়ে
দিয়েছে
দিয়েছেন
দিলেন
দু
দুটি
দুটো
দেয়
দেওয়া
দেওয়ার
দেখা
দেখে
দেখতে
দ্বারা
ধরে
ধরা
নয়
নানা
না
নাকি
নাগাদ
নিতে
নিজে
নিজেই
নিজের
নিজেদের
নিয়ে
নেওয়া
নেওয়ার
নেই
নাই
পক্ষে
পর্যন্ত
পাওয়া
পারেন
পারি
পারে
পরে
পরেই
পরেও
পর
পেয়ে
প্রতি
প্রভৃতি
প্রায়
ফের
ফলে
ফিরে
ব্যবহার
বলতে
বললেন
বলেছেন
বলল
বলা
বলেন
বলে
বহু
বসে
বার
বা
বিনা
বরং
বদলে
বাদে
বার
বিশেষ
বিভিন্ন
বিষয়টি
ব্যবহার
ব্যাপারে
ভাবে
ভাবেই
মধ্যে
মধ্যেই
মধ্যেও
মধ্যভাগে
মাধ্যমে
মাত্র
মতো
মতোই
মোটেই
যখন
যদি
যদিও
যাবে
যায়
যাকে
যাওয়া
যাওয়ার
যত
যতটা
যা
যার
যারা
যাঁর
যাঁরা
যাদের
যান
যাচ্ছে
যেতে
যাতে
যেন
যেমন
যেখানে
যিনি
যে
রেখে
রাখা
রয়েছে
রকম
শুধু
সঙ্গে
সঙ্গেও
সমস্ত
সব
সবার
সহ
সুতরাং
সহিত
সেই
সেটা
সেটি
সেটাই
সেটাও
সম্প্রতি
সেখান
সেখানে
সে
স্পষ্ট
স্বয়ং
হইতে
হইবে
হৈলে
হইয়া
হচ্ছে
হত
হতে
হতেই
হবে
হবেন
হয়েছিল
হয়েছে
হয়েছেন
হয়ে
হয়নি
হয়
হয়েই
হয়তো
হল
হলে
হলেই
হলেও
হলো
হিসাবে
হওয়া
হওয়ার
হওয়ায়
হন
হোক
জন
জনকে
জনের
জানতে
জানায়
জানিয়ে
জানানো
জানিয়েছে
জন্য
জন্যওজে
জে
বেশ
দেন
তুলে
ছিলেন
চান
চায়
চেয়ে
মোট
যথেষ্ট
টি
`,
	"ca": `
de
es
i
a
o
un
una
unes
uns
un
tot
també
altre
algun
alguna
alguns
algunes
ser
és
soc
ets
som
estic
està
estem
esteu
estan
com
en
per
perquè
per que
estat
estava
ans
abans
éssent
ambdós
però
per
poder
potser
puc
podem
podeu
poden
vaig
va
van
fer
faig
fa
fem
feu
fan
cada
fi
inclòs
primer
des de
conseguir
consegueixo
consigueix
consigueixes
conseguim
consigueixen
anar
haver
tenir
tinc
te
tenim
teniu
tene
el
la
les
els
seu
aquí
meu
teu
ells
elles
ens
nosaltres
vosaltres
si
dins
sols
solament
saber
saps
sap
sabem
sabeu
saben
últim
llarg
bastant
fas
molts
aquells
aquelles
seus
llavors
sota
dalt
ús
molt
era
eres
erem
eren
mode
bé
quant
quan
on
mentre
qui
amb
entre
sense
jo
aquell
`,
	"hr": `
a
ako
ali
bi
bih
bila
bili
bilo
bio
bismo
biste
biti
bumo
da
do
duž
ga
hoće
hoćemo
hoćete
hoćeš
hoću
i
iako
ih
ili
iz
ja
je
jedna
jedne
jedno
jer
jesam
jesi
jesmo
jest
jeste
jesu
jim
joj
još
ju
kada
kako
kao
koja
koje
koji
kojima
koju
kroz
li
me
mene
meni
mi
mimo
moj
moja
moje
mu
na
nad
nakon
nam
nama
nas
naš
naša
naše
našeg
ne
nego
neka
neki
nekog
neku
nema
netko
neće
nećemo
nećete
nećeš
neću
nešto
ni
nije
nikoga
nikoje
nikoju
nisam
nisi
nismo
niste
nisu
njega
njegov
njegova
njegovo
njemu
njezin
njezina
njezino
njih
njihov
njihova
njihovo
njim
njima
njoj
nju
no
o
od
odmah
on
ona
oni
ono
ova
pa
pak
po
pod
pored
prije
s
sa
sam
samo
se
sebe
sebi
si
smo
ste
su
sve
svi
svog
svoj
svoja
svoje
svom
ta
tada
taj
tako
te
tebe
tebi
ti
to
toj
tome
tu
tvoj
tvoja
tvoje
u
uz
vam
vama
vas
vaš
vaša
vaše
već
vi
vrlo
za
zar
će
ćemo
ćete
ćeš
ću
što
`,
	"cs": `
dnes
cz
timto
budes
budem
byli
jses
muj
svym
ta
tomto
tohle
tuto
tyto
jej
zda
proc
mate
tato
kam
tohoto
kdo
kteri
mi
nam
tom
tomuto
mit
nic
proto
kterou
byla
toho
protoze
asi
ho
nasi
napiste
re
coz
tim
takze
svych
jeji
svymi
jste
aj
tu
tedy
teto
bylo
kde
ke
prave
ji
nad
nejsou
ci
pod
tema
mezi
pres
ty
pak
vam
ani
kdyz
vsak
ne
jsem
tento
clanku
clanky
aby
jsme
pred
pta
jejich
byl
jeste
az
bez
take
pouze
prvni
vase
ktera
nas
novy
tipy
pokud
muze
design
strana
jeho
sve
jine
zpravy
nove
neni
vas
jen
podle
zde
clanek
uz
email
byt
vice
bude
jiz
nez
ktery
by
ktere
co
nebo
ten
tak
ma
pri
od
po
jsou
jak
dalsi
ale
si
ve
to
jako
za
zpet
ze
do
pro
je
na
`,
	"gl": `
a
aínda
alí
aquel
aquela
aquelas
aqueles
aquilo
aquí
ao
aos
as
así
á
ben
cando
che
co
coa
comigo
con
connosco
contigo
convosco
coas
cos
cun
cuns
cunha
cunhas
da
dalgunha
dalgunhas
dalgún
dalgúns
das
de
del
dela
delas
deles
desde
deste
do
dos
dun
duns
dunha
dunhas
e
el
ela
elas
eles
en
era
eran
esa
esas
ese
eses
esta
estar
estaba
está
están
este
estes
estiven
estou
eu
é
facer
foi
foron
fun
había
hai
iso
isto
la
las
lle
lles
lo
los
mais
me
meu
meus
min
miña
miñas
moi
na
nas
neste
nin
no
non
nos
nosa
nosas
noso
nosos
nós
nun
nunha
nuns
nunhas
o
os
ou
ó
ós
para
pero
pode
pois
pola
polas
polo
polos
por
que
se
senón
ser
seu
seus
sexa
sido
sobre
súa
súas
tamén
tan
te
ten
teñen
teño
ter
teu
teus
ti
tido
tiña
tiven
túa
túas
un
unha
unhas
uns
vos
vosa
vosas
voso
vosos
vós
`,
	"el": `
μή
ἑαυτοῦ
ἄν
ἀλλ’
ἀλλά
ἄλλοσ
ἀπό
ἄρα
αὐτόσ
δ’
δέ
δή
διά
δαί
δαίσ
ἔτι
ἐγώ
ἐκ
ἐμόσ
ἐν
ἐπί
εἰ
εἰμί
εἴμι
εἰσ
γάρ
γε
γα^
ἡ
ἤ
καί
κατά
μέν
μετά
μή
ὁ
ὅδε
ὅσ
ὅστισ
ὅτι
οὕτωσ
οὗτοσ
οὔτε
οὖν
οὐδείσ
οἱ
οὐ
οὐδέ
οὐκ
περί
πρόσ
σύ
σύν
τά
τε
τήν
τῆσ
τῇ
τι
τί
τισ
τίσ
τό
τοί
τοιοῦτοσ
τόν
τούσ
τοῦ
τῶν
τῷ
ὑμόσ
ὑπέρ
ὑπό
ὡσ
ὦ
ὥστε
ἐάν
παρά
σόσ
`,
	"he": `
את
לא
של
אני
על
זה
עם
כל
הוא
אם
או
גם
יותר
יש
לי
מה
אבל
פורום
אז
טוב
רק
כי
שלי
היה
בפורום
אין
עוד
היא
אחד
ב
ל
עד
לך
כמו
להיות
אתה
כמה
אנחנו
הם
כבר
אנשים
אפשר
תודה
שלא
אותו
ה
מאוד
הרבה
ולא
ממש
לו
א
מי
חיים
בית
שאני
יכול
שהוא
כך
הזה
איך
היום
קצת
עכשיו
שם
בכל
יהיה
תמיד
י
שלך
הכי
ש
בו
לעשות
צריך
כן
פעם
לכם
ואני
משהו
אל
שלו
שיש
ו
וגם
אתכם
אחרי
בנושא
כדי
פשוט
לפני
שזה
אותי
אנו
למה
דבר
כ
כאן
אולי
טובים
רוצה
שנה
בעלי
החיים
למען
אתם
מ
בין
יום
זאת
איזה
ביותר
לה
אחת
הכל
הפורומים
לכל
אלא
פה
יודע
שלום
דקות
לנו
השנה
דרך
אדם
נראה
זו
היחידה
רוצים
בכלל
טובה
שלנו
האם
הייתי
הלב
היו
ח
שדרות
בלי
להם
שאתה
אותה
מקום
ואתם
חלק
בן
בואו
אחר
האחת
אותך
כמובן
בגלל
באמת
מישהו
ילדים
אותם
הפורום
טיפוח
וזה
ר
שהם
אך
מזמין
ישראל
כוס
זמן
ועוד
הילדים
עדיין
כזה
עושה
שום
לקחת
העולם
תפוז
לראות
לפורום
וכל
לקבל
נכון
יוצא
לעולם
גדול
אפילו
ניתן
שני
אוכל
קשה
משחק
ביום
ככה
אמא
בת
השבוע
נוספים
לגבי
בבית
אחרת
לפי
ללא
שנים
הזמן
שמן
מעט
לפחות
אף
שוב
שלהם
במקום
כולם
נועית
הבא
מעל
לב
המון
לדבר
ע
אוהב
מוסיפים
חצי
בעיקר
כפות
לפעמים
שהיא
הנהלת
ועל
ק
אוהבים
ת
יודעת
ד
גרוע
שאנחנו
מים
לילדים
בארץ
מודיע
אשמח
שלכם
פחות
לכולם
די
אהבה
יכולה
דברים
הקהילה
לעזור
פרטים
בדיוק
מלח
קל
הראשי
שלה
להוסיף
השני
לדעתי
בר
למרות
שגם
מוזמנים
לאחר
במה
חושב
מאד
יפה
להגשים
חדש
קטן
מחפשים
שמח
מדברים
ואם
במיוחד
עבודה
מדי
ואז
חשוב
שאם
אוהבת
פעמים
מנהלת
אומר
מול
קשר
מנהל
שיהיה
שאין
שאנו
האהבה
ס
הצטרפו
כפית
בשביל
החגים
אופן
לתת
כף
בתוך
סוכר
גיל
בהצלחה
והוא
מקווה
סתם
ויש
נגד
כמעט
שאת
עולה
אי
מספר
ראשון
לדרך
נהיה
לעצב
עושים
ולנהל
היתה
עליו
מזה
הייתה
בא
בפרס
חלות
ראש
מזמינים
טיפים
מכבי
רבה
הורים
‡
מקרה
קרן
המוצלח
להגיע
גדולה
כנראה
החמשיר
הראשון
פלפל
המשחק
וכאן
לדעת
ואת
גרועים
ספר
אגב
אחרים
להגיד
בתפוז
והעולם
אופנה
דווקא
מספיק
שעות
תמונות
כשאנחנו
שוקולד
ולכן
ג
לקרוא
לניהול
שבוע
ויופי
חלום
בה
שהיה
שאלה
מקומה
הזו
בפורומים
החדש
מתאמצים
שחקן
שמזינים
נשמת
בערך
מכל
ומה
רגל
כסף
רואה
קטנה
בצל
בעולם
אינטרנט
חוץ
ברור
הולך
חושבת
לזה
כלום
הן
כאלה
בטוח
הדבר
תהיה
מגיע
סוף
האמת
ממנו
מיכל
החדשה
לתרום
האנשים
ועד
בדרך
אצלי
ההורים
בני
מתוך
כאשר
לבד
ראיתי
מצב
מלא
לבחור
נשמח
החג
רע
עוף
מן
להביא
מצאתי
כתובות
מעניין
צריכה
להכנס
לחלוטין
שעה
מתכון
קודם
תשובות
מדובר
ניהול
מזל
כדאי
יהיו
ההודעות
בוקר
נילוות
איפה
בעיה
קמח
ללכת
פורומים
אמר
נושא
ההכנה
בבקשה
שכל
הזאת
למשחק
פנינה
תחרות
חבר
לקנות
מהם
רגע
גרם
אלו
עצמו
מראש
הכלב
כולנו
עדיף
איתו
למשל
לבשל
למי
רעיונות
הבלוג
רוב
אביב
כרגע
בסוף
אלה
לחג
ערוץ
שווה
באופן
מאמין
לבן
בזה
הכבוד
לראש
ם
ימי
שחור
בצורה
בעמוד
ועם
וחצי
האלה
תמונה
בשלב
משחקים
נו
`,
	"hi": `
के
का
एक
में
की
है
यह
और
से
हैं
को
पर
इस
होता
कि
जो
कर
मे
गया
करने
किया
लिये
अपने
ने
बनी
नहीं
तो
ही
या
एवं
दिया
हो
इसका
था
द्वारा
हुआ
तक
साथ
करना
वाले
बाद
लिए
आप
कुछ
सकते
किसी
ये
इसके
सबसे
इसमें
थे
दो
होने
वह
वे
करते
बहुत
कहा
वर्ग
कई
करें
होती
अपनी
उनके
थी
यदि
हुई
जा
ना
इसे
कहते
जब
होते
कोई
हुए
व
न
अभी
जैसे
सभी
करता
उनकी
तरह
उस
आदि
कुल
एस
रहा
इसकी
सकता
रहे
उनका
इसी
रखें
अपना
पे
उसके
`,
	"ga": `
a
ach
ag
agus
an
aon
ar
arna
as
b'
ba
beirt
bhúr
caoga
ceathair
ceathrar
chomh
chtó
chuig
chun
cois
céad
cúig
cúigear
d'
daichead
dar
de
deich
deichniúr
den
dhá
do
don
dtí
dá
dár
dó
faoi
faoin
faoina
faoinár
fara
fiche
gach
gan
go
gur
haon
hocht
i
iad
idir
in
ina
ins
inár
is
le
leis
lena
lenár
m'
mar
mo
mé
na
nach
naoi
naonúr
ná
ní
níor
nó
nócha
ocht
ochtar
os
roimh
sa
seacht
seachtar
seachtó
seasca
seisear
siad
sibh
sinn
sna
sé
sí
tar
thar
thú
triúr
trí
trína
trínár
tríocha
tú
um
ár
é
éis
í
ó
ón
óna
ónár
`,
	"ja": `
これ
それ
あれ
この
その
あの
ここ
そこ
あそこ
こちら
どこ
だれ
なに
なん
何
私
貴方
貴方方
我々
私達
あの人
あのかた
彼女
彼
です
あります
おります
います
は
が
の
に
を
で
え
から
まで
より
も
どの
と
し
それで
しかし
`,
	"ku": `
# and
و
# which
کە
# of
ی
# made/did
کرد
# that/which
ئەوەی
# on/head
سەر
# two
دوو
# also
هەروەها
# from/that
لەو
# makes/does
دەکات
# some
چەند
# every
هەر

# demonstratives
# that
ئەو
# this
ئەم

# personal pronouns
# I
من
# we
ئێمە
# you
تۆ
# you
ئێوە
# he/she/it
ئەو
# they
ئەوان

# prepositions
# to/with/by
بە
پێ
# without
بەبێ
# along with/while/during
بەدەم
# in the opinion of
بەلای
# according to
بەپێی
# before
بەرلە
# in the direction of
بەرەوی
# in front of/toward
بەرەوە
# before/in the face of
بەردەم
# without
بێ
# except for
بێجگە
# for
بۆ
# on/in
دە
تێ
# with
دەگەڵ
# after
دوای
# except for/aside from
جگە
# in/from
لە
لێ
# in front of/before/because of
لەبەر
# between/among
لەبەینی
# concerning/about
لەبابەت
# concerning
لەبارەی
# instead of
لەباتی
# beside
لەبن
# instead of
لەبرێتی
# behind
لەدەم
# with/together with
لەگەڵ
# by
لەلایەن
# within
لەناو
# between/among
لەنێو
# for the sake of
لەپێناوی
# with respect to
لەرەوی
# by means of/for
لەرێ
# for the sake of
لەرێگا
# on/on top of/according to
لەسەر
# under
لەژێر
# between/among
ناو
# between/among
نێوان
# after
پاش
# before
پێش
# like
وەک
`,
	"lv": `
ārpus
šaipus
aiz
ap
apakš
apakšpus
arī
ar
ar
augšpus
būšu     
būs
būsi
būsiet
būsim
būt  
bet
bez
bijām
bijāt
bija
biji
biju
caur
dēļ
diemžēl
diezin
droši
esam
esat
esi
esmu
gan
gar
iekām
iekāms
iekš
iekšpus
iekam
iekams
ik
ir
it
itin
iz
jā
ja
jau
jebšu
jeb
jel
jo
kā
kļūšu
kļūs
kļūsi
kļūsiet
kļūsim
kļūst
kļūstam
kļūstat
kļūsti
kļūstu
kļūt
kļuvām
kļuvāt
kļuva
kļuvi
kļuvu
ka
kamēr
kaut
kolīdz
kopš
līdz
līdzko
labad
lai
lejpus
nē
ne
nebūt
nedz
nekā
nevis
nezin
no
nu
otrpus
pār
pēc
pa
par
pat
pie
pirms
pret
priekš
starp
tā
tādēļ
tālab
tāpēc
taču
tad
tak
tapāt
tapšu
tapi
taps
tapsi
tapsiet
tapsim
tapt
te
tiec
tiek
tiekam
tiekat
tieku
tikām
tikāt
tikšu
tik
tika
tikai
tiki
tikko
tiklīdz
tiklab
tiks
tiksiet
tiksim
tikt
tiku
tikvien
tomēr
topat
turpretī
turpretim
un
uz
vai
varēšu
varējām
varējāt
varēja
varēji
varēju
varēs
varēsi
varēsiet
varēsim
varēt
var
varat
viņpus
vien
vien
virs
virspus
vis
zem
`,
	"lt": `
á
ákypai
ástriþai
ðájá
ðalia
ðe
ðiàjà
ðiàja
ðiàsias
ðiøjø
ðiøjø
ði
ðiaisiais
ðiajai
ðiajam
ðiajame
ðiapus
ðiedvi
ðieji
ðiesiems
ðioji
ðiojo
ðiojoje
ðiokia
ðioks
ðiosiomis
ðiosioms
ðiosios
ðiosios
ðiosiose
ðis
ðisai
ðit
ðita
ðitas
ðitiedvi
ðitokia
ðitoks
ðituodu
ðiuodu
ðiuoju
ðiuosiuose
ðiuosius
ðtai
þemiau
að
abi
abidvi
abiejø
abiejose
abiejuose
abiem
abigaliai
abipus
abu
abudu
ai
anàjà
anàjá
anàja
anàsias
anøjø
anøjø
ana
ana
anaiptol
anaisiais
anajai
anajam
anajame
anapus
anas
anasai
anasis
anei
aniedvi
anieji
aniesiems
anoji
anojo
anojoje
anokia
anoks
anosiomis
anosioms
anosios
anosios
anosiose
anot
ant
antai
anuodu
anuoju
anuosiuose
anuosius
apie
aplink
ar
ar
arba
argi
arti
aukðèiau
be
be
bei
beje
bemaþ
bent
bet
betgi
beveik
dëka
dël
dëlei
dëlto
dar
dargi
daugmaþ
deja
ech
et
gal
galbût
galgi
gan
gana
gi
greta
ið
iðilgai
iðvis
idant
iki
iki
ir
irgi
it
itin
jàjà
jàja
jàsias
jájá
jøjø
jøjø
jûsø
jûs
jûsiðkë
jûsiðkis
jaisiais
jajai
jajam
jajame
jei
jeigu
ji
jiedu
jiedvi
jieji
jiesiems
jinai
jis
jisai
jog
joji
jojo
jojoje
jokia
joks
josiomis
josioms
josios
josios
josiose
judu
judvi
juk
jumis
jums
jumyse
juodu
juoju
juosiuose
juosius
jus
kaþin
kaþkas
kaþkatra
kaþkatras
kaþkokia
kaþkoks
kaþkuri
kaþkuris
kad
kada
kadangi
kai
kaip
kaip
kaipgi
kas
katra
katras
katriedvi
katruodu
kiaurai
kiek
kiekvienas
kieno
kita
kitas
kitokia
kitoks
kodël
kokia
koks
kol
kolei
kone
kuomet
kur
kurgi
kuri
kuriedvi
kuris
kuriuodu
lai
lig
ligi
link
lyg
mûsø
mûsiðkë
mûsiðkis
maþdaug
maþne
manàjà
manàjá
manàja
manàsias
manæs
manøjø
manøjø
man
manaisiais
manajai
manajam
manajame
manas
manasai
manasis
mane
maniðkë
maniðkis
manieji
maniesiems
manim
manimi
mano
manoji
manojo
manojoje
manosiomis
manosioms
manosios
manosios
manosiose
manuoju
manuosiuose
manuosius
manyje
mat
mes
mudu
mudvi
mumis
mums
mumyse
mus
në
na
nagi
ne
nebe
nebent
nebent
negi
negu
nei
nei
nejau
nejaugi
nekaip
nelyginant
nes
net
netgi
netoli
neva
nors
nors
nuo
o
ogi
ogi
oi
paèiø
paèiais
paèiam
paèiame
paèiu
paèiuose
paèius
paeiliui
pagal
pakeliui
palaipsniui
palei
pas
pasak
paskos
paskui
paskum
patá
pat
pati
patiems
paties
pats
patys
per
per
pernelyg
pirm
pirma
pirmiau
po
prieð
prieðais
prie
pro
pusiau
rasi
rodos
sau
savàjà
savàjá
savàja
savàsias
savæs
savøjø
savøjø
savaisiais
savajai
savajam
savajame
savas
savasai
savasis
save
saviðkë
saviðkis
savieji
saviesiems
savimi
savo
savoji
savojo
savojoje
savosiomis
savosioms
savosios
savosios
savosiose
savuoju
savuosiuose
savuosius
savyje
skersai
skradþiai
staèiai
su
sulig
tàjà
tàjá
tàja
tàsias
tøjø
tøjø
tûlas
taèiau
ta
tad
tai
tai
taigi
taigi
taip
taipogi
taisiais
tajai
tajam
tajame
tamsta
tarp
tarsi
tarsi
tartum
tartum
tarytum
tas
tasai
tau
tavàjà
tavàjá
tavàja
tavàsias
tavæs
tavøjø
tavøjø
tavaisiais
tavajai
tavajam
tavajame
tavas
tavasai
tavasis
tave
taviðkë
taviðkis
tavieji
taviesiems
tavimi
tavo
tavoji
tavojo
tavojoje
tavosiomis
tavosioms
tavosios
tavosios
tavosiose
tavuoju
tavuosiuose
tavuosius
tavyje
te
tegu
tegu
tegul
tegul
tiedvi
tieji
ties
tiesiems
tiesiog
tik
tik
tikriausiai
tiktai
tiktai
toji
tojo
tojoje
tokia
toks
tol
tolei
toliau
tosiomis
tosioms
tosios
tosios
tosiose
tu
tuodu
tuoju
tuosiuose
tuosius
turbût
uþ
uþtat
uþvis
vël
vëlgi
va
vai
viduj
vidury
vien
vienas
vienokia
vienoks
vietoj
virð
virðuj
virðum
vis
vis dëlto
visa
visas
visgi
visokia
visoks
vos
vos
ypaè
`,
	"mr": `
आहे
या
आणि
व
नाही  
आहेत
यानी
हे
तर
ते
असे
होते
केली
हा
ही
पण
करणयात
काही
केले
एक
केला
अशी
मात्र  
त्यानी
सुरू
करून
होती
असून
आले
त्यामुळे
झाली
होता
दोन
झाले
मुबी
होत
त्या
आता
असा
याच्या
त्याच्या
ता
आली
की
पम
तो
झाला
त्री
तरी
म्हणून
त्याना
अनेक
काम
माहिती
हजार
सागित्ले
दिली
आला
आज
ती
तसेच
एका
याची
येथील
सर्व
न
डॉ
तीन
येथे
पाटील
असलयाचे
त्याची
काय
आपल्या
म्हणजे
याना
म्हणाले
त्याचा
असलेल्या
मी
गेल्या
याचा
येत
म
लाख
कमी
जात    
टा
होणार
किवा
का
अधिक
घेऊन      
परयतन
कोटी
झालेल्या
निर्ण्य
येणार
व्यकत
`,
	"fa": `
و
در
به
از
كه
مي
اين
است
را
با
هاي
براي
آن
يك
شود
شده
خود
ها
كرد
شد
اي
تا
كند
بر
بود
گفت
نيز
وي
هم
كنند
دارد
ما
كرده
يا
اما
بايد
دو
اند
هر
خواهد
او
مورد
آنها
باشد
ديگر
مردم
نمي
بين
پيش
پس
اگر
همه
صورت
يكي
هستند
بي
من
دهد
هزار
نيست
استفاده
داد
داشته
راه
داشت
چه
همچنين
كردند
داده
بوده
دارند
همين
ميليون
سوي
شوند
بيشتر
بسيار
روي
گرفته
هايي
تواند
اول
نام
هيچ
چند
جديد
بيش
شدن
كردن
كنيم
نشان
حتي
اينكه
ولی
توسط
چنين
برخي
نه
ديروز
دوم
درباره
بعد
مختلف
گيرد
شما
گفته
آنان
بار
طور
گرفت
دهند
گذاري
بسياري
طي
بودند
ميليارد
بدون
تمام
كل
تر
براساس
شدند
ترين
امروز
باشند
ندارد
چون
قابل
گويد
ديگري
همان
خواهند
قبل
آمده
اكنون
تحت
طريق
گيري
جاي
هنوز
چرا
البته
كنيد
سازي
سوم
كنم
بلكه
زير
توانند
ضمن
فقط
بودن
حق
آيد
وقتي
اش
يابد
نخستين
مقابل
خدمات
امسال
تاكنون
مانند
تازه
آورد
فكر
آنچه
نخست
نشده
شايد
چهار
جريان
پنج
ساخته
زيرا
نزديك
برداري
كسي
ريزي
رفت
گردد
مثل
آمد
ام
بهترين
دانست
كمتر
دادن
تمامي
جلوگيري
بيشتري
ايم
ناشي
چيزي
آنكه
بالا
بنابراين
ايشان
بعضي
دادند
داشتند
برخوردار
نخواهد
هنگام
نبايد
غير
نبود
ديده
وگو
داريم
چگونه
بندي
خواست
فوق
ده
نوعي
هستيم
ديگران
همچنان
سراسر
ندارند
گروهي
سعي
روزهاي
آنجا
يكديگر
كردم
بيست
بروز
سپس
رفته
آورده
نمايد
باشيم
گويند
زياد
خويش
همواره
گذاشته
شش
نداشته
شناسي
خواهيم
آباد
داشتن
نظير
همچون
باره
نكرده
شان
سابق
هفت
دانند
جايي
بی
جز
زیرِ
رویِ
سریِ
تویِ
جلویِ
پیشِ
عقبِ
بالایِ
خارجِ
وسطِ
بیرونِ
سویِ
کنارِ
پاعینِ
نزدِ
نزدیکِ
دنبالِ
حدودِ
برابرِ
طبقِ
مانندِ
ضدِّ
هنگامِ
برایِ
مثلِ
بارة
اثرِ
تولِ
علّتِ
سمتِ
عنوانِ
قصدِ
روب
جدا
کی
که
چیست
هست
کجا
کجاست
کَی
چطور
کدام
آیا
مگر
چندین
یک
چیزی
دیگر
کسی
بعری
هیچ
چیز
جا
کس
هرگز
یا
تنها
بلکه
خیاه
بله
بلی
آره
آری
مرسی
البتّه
لطفاً
ّه
انکه
وقتیکه
همین
پیش
مدّتی
هنگامی
مان
تان
`,
	"ro": `
acea
aceasta
această
aceea
acei
aceia
acel
acela
acele
acelea
acest
acesta
aceste
acestea
aceşti
aceştia
acolo
acord
acum
ai
aia
aibă
aici
al
ăla
ale
alea
ălea
altceva
altcineva
am
ar
are
aş
aşadar
asemenea
asta
ăsta
astăzi
astea
ăstea
ăştia
asupra
aţi
au
avea
avem
aveţi
azi
bine
bucur
bună
ca
că
căci
când
care
cărei
căror
cărui
cât
câte
câţi
către
câtva
caut
ce
cel
ceva
chiar
cinci
cînd
cine
cineva
cît
cîte
cîţi
cîtva
contra
cu
cum
cumva
curând
curînd
da
dă
dacă
dar
dată
datorită
dau
de
deci
deja
deoarece
departe
deşi
din
dinaintea
dintr-
dintre
doi
doilea
două
drept
după
ea
ei
el
ele
eram
este
eşti
eu
face
fără
fata
fi
fie
fiecare
fii
fim
fiţi
fiu
frumos
graţie
halbă
iar
ieri
îi
îl
îmi
împotriva
în
înainte
înaintea
încât
încît
încotro
între
întrucât
întrucît
îţi
la
lângă
le
li
lîngă
lor
lui
mă
mai
mâine
mea
mei
mele
mereu
meu
mi
mie
mîine
mine
mult
multă
mulţi
mulţumesc
ne
nevoie
nicăieri
nici
nimeni
nimeri
nimic
nişte
noastră
noastre
noi
noroc
noştri
nostru
nouă
nu
opt
ori
oricând
oricare
oricât
orice
oricînd
oricine
oricît
oricum
oriunde
până
patra
patru
patrulea
pe
pentru
peste
pic
pînă
poate
pot
prea
prima
primul
prin
puţin
puţina
puţină
rog
sa
să
săi
sale
şapte
şase
sau
său
se
şi
sînt
sîntem
sînteţi
spate
spre
ştiu
sub
sunt
suntem
sunteţi
sută
ta
tăi
tale
tău
te
ţi
ţie
timp
tine
toată
toate
tot
toţi
totuşi
trei
treia
treilea
tu
un
una
unde
undeva
unei
uneia
unele
uneori
unii
unor
unora
unu
unui
unuia
unul
vă
vi
voastră
voastre
voi
voştri
vostru
vouă
vreme
vreo
vreun
zece
zero
zi
zice
`,
	"sk": `
a
aby
aj
ak
ako
ale
alebo
and
ani
áno
asi
až
bez
bude
budem
budeš
budeme
budete
budú
by
bol
bola
boli
bolo
byť
cez
čo
či
ďalší
ďalšia
ďalšie
dnes
do
ho
ešte
for
i
ja
je
jeho
jej
ich
iba
iné
iný
som
si
sme
sú
k
kam
každý
každá
každé
každí
kde
keď
kto
ktorá
ktoré
ktorou
ktorý
ktorí
ku
lebo
len
ma
mať
má
máte
medzi
mi
mna
mne
mnou
musieť
môcť
môj
môže
my
na
nad
nám
náš
naši
nie
nech
než
nič
niektorý
nové
nový
nová
nové
noví
o
od
odo
of
on
ona
ono
oni
ony
po
pod
podľa
pokiaľ
potom
práve
pre
prečo
preto
pretože
prvý
prvá
prvé
prví
pred
predo
pri
pýta
s
sa
so
si
svoje
svoj
svojich
svojím
svojími
ta
tak
takže
táto
teda
te
tě
ten
tento
the
tieto
tým
týmto
tiež
to
toto
toho
tohoto
tom
tomto
tomuto
toto
tu
tú
túto
tvoj
ty
tvojími
už
v
vám
váš
vaše
vo
viac
však
všetok
vy
z
za
zo
že
`,
	"th": `
ไว้
ไม่
ไป
ได้
ให้
ใน
โดย
แห่ง
แล้ว
และ
แรก
แบบ
แต่
เอง
เห็น
เลย
เริ่ม
เรา
เมื่อ
เพื่อ
เพราะ
เป็นการ
เป็น
เปิดเผย
เปิด
เนื่องจาก
เดียวกัน
เดียว
เช่น
เฉพาะ
เคย
เข้า
เขา
อีก
อาจ
อะไร
ออก
อย่าง
อยู่
อยาก
หาก
หลาย
หลังจาก
หลัง
หรือ
หนึ่ง
ส่วน
ส่ง
สุด
สําหรับ
ว่า
วัน
ลง
ร่วม
ราย
รับ
ระหว่าง
รวม
ยัง
มี
มาก
มา
พร้อม
พบ
ผ่าน
ผล
บาง
น่า
นี้
นํา
นั้น
นัก
นอกจาก
ทุก
ที่สุด
ที่
ทําให้
ทํา
ทาง
ทั้งนี้
ง
ถ้า
ถูก
ถึง
ต้อง
ต่างๆ
ต่าง
ต่อ
ตาม
ตั้งแต่
ตั้ง
ด้าน
ด้วย
ดัง
ซึ่ง
ช่วง
จึง
จาก
จัด
จะ
คือ
ความ
ครั้ง
คง
ขึ้น
ของ
ขอ
ขณะ
ก่อน
ก็
การ
กับ
กัน
กว่า
กล่าว
`,
	"tr": `
acaba
altmýþ
altý
ama
bana
bazý
belki
ben
benden
beni
benim
beþ
bin
bir
biri
birkaç
birkez
birþey
birþeyi
biz
bizden
bizi
bizim
bu
buna
bunda
bundan
bunu
bunun
da
daha
dahi
de
defa
diye
doksan
dokuz
dört
elli
en
gibi
hem
hep
hepsi
her
hiç
iki
ile
INSERmi
ise
için
katrilyon
kez
ki
kim
kimden
kime
kimi
kýrk
milyar
milyon
mu
mü
mý
nasýl
ne
neden
nerde
nerede
nereye
niye
niçin
on
ona
ondan
onlar
onlardan
onlari
onlarýn
onu
otuz
sanki
sekiz
seksen
sen
senden
seni
senin
siz
sizden
sizi
sizin
trilyon
tüm
ve
veya
ya
yani
yedi
yetmiþ
yirmi
yüz
çok
çünkü
üç
þey
þeyden
þeyi
þeyler
þu
þuna
þunda
þundan
þunu
`,
	"uk": `
з
й
що
та
Із
але
цей
коли
як
чого
хоча
нам
яко╞
чи
це
в╡д
про
╡
╞х
╙
Інших
ти
вІн
вона
воно
ми
ви
вони
`,
	"ur": `
ثھی
خو
گی
اپٌے
گئے
ثہت
طرف
ہوبری
پبئے
اپٌب
دوضری
گیب
کت
گب
ثھی
ضے
ہر
پر
اش
دی
گے
لگیں
ہے
ثعذ
ضکتے
وٍ
تھی
اى
دیب
لئے
والے
یہ
ثدبئے
ضکتی
ًے
تھب
اًذر
رریعے
لگی
ہوبرا
ہوًے
ثبہر
ضکتب
ًہیں
تو
اور
رہب
لگے
ہوضکتب
ہوں
کب
ہوبرے
توبم
کیب
ایطے
رہی
هگر
ہوضکتی
ہیں
کریں
ہو
تک
کی
ایک
رہے
هیں
ہوضکتے
کیطے
ہوًب
تت
کہ
ہوا
آئے
ضبتھ
ًب
تھے
کیوں
ہوتب
ًہ
خت
کے
پھر
ثغیر
خبرہے
رکھ
کیطب
کوئی
رریعے
ثبرے
خب
اضطرذ
ثلکہ
خجکہ
رکھتب
کیطرف
ثراں
خبرہب
رریعہ
کطی
اضکب
ثٌذ
خص
رکھتبہوں
کیلئے
ثبئیں
توہیں
دوضرے
کررہی
اضکی
ثیچ
خوکہ
رکھتی
کیوًکہ
دوًوں
کررہے
خبرہی
ثرآں
اضکے
پچھلا
خیطب
رکھتے
کےثعذ
توہی
دوراى
کررہب
یہبں
ٓ آش
اًہیں
ثي
پطٌذ
تھوڑا
چکے
زکویہ
دوضروں
ضکب
اة
اوًچب
ثٌب
پل
تھوڑی
چلا
خبهوظ
دیتب
ضکٌب
اخبزت
اوًچبئی
ثٌبرہب
پوچھب
تھوڑے
چلو
ختن
دیتی
ضکی
اچھب
اوًچی
ثٌبرہی
پوچھتب
تیي
چلیں
در
دیتے
ضکے
اچھی
اوًچے
ثٌبرہے
پوچھتی
خبًب
چلے
درخبت
دیر
ضلطلہ
اچھے
اٹھبًب
ثٌبًب
پوچھتے
خبًتب
چھوٹب
درخہ
دیکھٌب
ضوچ
اختتبم
اہن
ثٌذ
پوچھٌب
خبًتی
چھوٹوں
درخے
دیکھو
ضوچب
ادھر
آئی
ثٌذکرًب
پوچھو
خبًتے
چھوٹی
درزقیقت
دیکھی
ضوچتب
ارد
آئے
ثٌذکرو
پوچھوں
خبًٌب
چھوٹے
درضت
دیکھیں
ضوچتی
اردگرد
آج
ثٌذی
پوچھیں
خططرذ
چھہ
دش
دیٌب
ضوچتے
ارکبى
آخر
ثڑا
پورا
خگہ
چیسیں
دفعہ
دے
ضوچٌب
اضتعوبل
آخرکبر
ثڑوں
پہلا
خگہوں
زبصل
دکھبئیں
راضتوں
ضوچو
اضتعوبلات
آدهی
ثڑی
پہلی
خگہیں
زبضر
دکھبتب
راضتہ
ضوچی
اغیب
آًب
ثڑے
پہلےضی
خلذی
زبل
دکھبتی
راضتے
ضوچیں
اطراف
آٹھ
ثھر
پہلےضے
خٌبة
زبل
دکھبتے
رکي
ضیذھب
افراد
آیب
ثھرا
پہلےضےہی
خواى
زبلات
دکھبًب
رکھب
ضیذھی
اکثر
ثب
ثھراہوا
پیع
خوًہی
زبلیہ
دکھبو
رکھی
ضیذھے
اکٹھب
ثبترتیت
ثھرپور
تبزٍ
خیطبکہ
زصوں
دکھبیب
رکھے
ضیکٌڈ
اکٹھی
ثبری
ثہتر
تر
چبر
زصہ
دلچطپ
زیبدٍ
غبیذ
اکٹھے
ثبلا
ثہتری
ترتیت
چبہب
زصے
دلچطپی
ضبت
غخص
اکیلا
ثبلترتیت
ثہتریي
تریي
چبہٌب
زقبئق
دلچطپیبں
ضبدٍ
غذ
اکیلی
ثرش
پبش
تعذاد
چبہے
زقیتیں
هٌبضت
ضبرا
غروع
اکیلے
ثغیر
پبًب
اً
ب
ج
ی
ر
ق
ت
چکب
زقیقت
دو
ضبرے
غروعبت
اگرچہ
ثلٌذ
پبًچ
تن
چکی
زکن
دور
ضبل
غے
الگ
ثلٌذوثبلا
پراًب
تٌہب
چکیں
اً
ب
و
ک
ز
دوضرا
ضبلوں
صبف
صسیر
قجیلہ
کوًطے
لازهی
هطئلے
ًیب
طریق
کرتی
کہتے
صفر
قطن
کھولا
لگتب
هطبئل
وار
طریقوں
کرتے
کہٌب
صورت
کئی
کھولٌب
لگتی
هطتعول
وار
طریقہ
کرتےہو
کہٌب
صورتسبل
کئے
کھولو
لگتے
هػتول
ٹھیک
طریقے
کرًب
کہو
صورتوں
کبفی
کھولی
لگٌب
هطلق
ڈھوًڈا
طور
کرو
کہوں
صورتیں
کبم
کھولیں
لگی
هعلوم
ڈھوًڈلیب
طورپر
کریں
کہی
ضرور
کجھی
کھولے
لگے
هکول
ڈھوًڈًب
ظبہر
کرے
کہیں
ضرورت
کرا
کہب
لوجب
هلا
ڈھوًڈو
عذد
کل
کہیں
اً
ب
ت
ر
و
ر
ض
کرتب
کہتب
لوجی
هوکي
ڈھوًڈی
عظین
کن
کہے
ضروری
کرتبہوں
کہتی
لوجے
هوکٌبت
ڈھوًڈیں
علاقوں
کوتر
کیے

لوسبت
هوکٌہ
ہن
لے
ًبپطٌذ
ہورہے
علاقہ
کورا
کےرریعے
لوسہ
هڑا
ہوئی
هتعلق
ًبگسیر
ہوگئی
علاقے
کوروں
گئی
لو
هڑًب
ہوئے
هسترم
ًطجت
ہوگئے
علاوٍ
کورٍ
گرد
لوگ
هڑے
ہوتی
هسترهہ
ًقطہ
ہوگیب
اً
ب
ه
و
و
ع
کورے
گروپ
لوگوں
هہرثبى
ہوتے
هسطوش
ًکبلٌب
ہوًی
عووهی
کوطي
گروٍ
لڑکپي
هیرا
ہوچکب
هختلف
ًکتہ
ہی
فرد
کوى
گروہوں
لی
هیری
ہوچکی
هسیذ
ًو
اً
ب
ی
ق
ی
فی
کوًطب
گٌتی
لیب
هیرے
ہوچکے
هطئلہ
ًوخواى
یقیٌی
قجل
کوًطی
اً
ب
ه
ز
لا
لیٌب
ًئی
ہورہب
لیں
ًئے
ہورہی
ثبعث
ضت
`,
	"vi": `
a lô
a ha
ai
ai ai
ai nấy
ai đó
alô
amen
anh
anh ấy
ba
ba ba
ba bản
ba cùng
ba họ
ba ngày
ba ngôi
ba tăng
bao giờ
bao lâu
bao nhiêu
bao nả
bay biến
biết
biết bao
biết bao nhiêu
biết chắc
biết chừng nào
biết mình
biết mấy
biết thế
biết trước
biết việc
biết đâu
biết đâu chừng
biết đâu đấy
biết được
buổi
buổi làm
buổi mới
buổi ngày
buổi sớm
bà
bà ấy
bài
bài bác
bài bỏ
bài cái
bác
bán
bán cấp
bán dạ
bán thế
bây bẩy
bây chừ
bây giờ
bây nhiêu
bèn
béng
bên
bên bị
bên có
bên cạnh
bông
bước
bước khỏi
bước tới
bước đi
bạn
bản
bản bộ
bản riêng
bản thân
bản ý
bất chợt
bất cứ
bất giác
bất kì
bất kể
bất kỳ
bất luận
bất ngờ
bất nhược
bất quá
bất quá chỉ
bất thình lình
bất tử
bất đồ
bấy
bấy chầy
bấy chừ
bấy giờ
bấy lâu
bấy lâu nay
bấy nay
bấy nhiêu
bập bà bập bõm
bập bõm
bắt đầu
bắt đầu từ
bằng
bằng cứ
bằng không
bằng người
bằng nhau
bằng như
bằng nào
bằng nấy
bằng vào
bằng được
bằng ấy
bển
bệt
bị
bị chú
bị vì
bỏ
bỏ bà
bỏ cha
bỏ cuộc
bỏ không
bỏ lại
bỏ mình
bỏ mất
bỏ mẹ
bỏ nhỏ
bỏ quá
bỏ ra
bỏ riêng
bỏ việc
bỏ xa
bỗng
bỗng chốc
bỗng dưng
bỗng không
bỗng nhiên
bỗng nhưng
bỗng thấy
bỗng đâu
bộ
bộ thuộc
bộ điều
bội phần
bớ
bởi
bởi ai
bởi chưng
bởi nhưng
bởi sao
bởi thế
bởi thế cho nên
bởi tại
bởi vì
bởi vậy
bởi đâu
bức
cao
cao lâu
cao ráo
cao răng
cao sang
cao số
cao thấp
cao thế
cao xa
cha
cha chả
chao ôi
chia sẻ
chiếc
cho
cho biết
cho chắc
cho hay
cho nhau
cho nên
cho rằng
cho rồi
cho thấy
cho tin
cho tới
cho tới khi
cho về
cho ăn
cho đang
cho được
cho đến
cho đến khi
cho đến nỗi
choa
chu cha
chui cha
chung
chung cho
chung chung
chung cuộc
chung cục
chung nhau
chung qui
chung quy
chung quy lại
chung ái
chuyển
chuyển tự
chuyển đạt
chuyện
chuẩn bị
chành chạnh
chí chết
chính
chính bản
chính giữa
chính là
chính thị
chính điểm
chùn chùn
chùn chũn
chú
chú dẫn
chú khách
chú mày
chú mình
chúng
chúng mình
chúng ta
chúng tôi
chúng ông
chăn chắn
chăng
chăng chắc
chăng nữa
chơi
chơi họ
chưa
chưa bao giờ
chưa chắc
chưa có
chưa cần
chưa dùng
chưa dễ
chưa kể
chưa tính
chưa từng
chầm chập
chậc
chắc
chắc chắn
chắc dạ
chắc hẳn
chắc lòng
chắc người
chắc vào
chắc ăn
chẳng lẽ
chẳng những
chẳng nữa
chẳng phải
chết nỗi
chết thật
chết tiệt
chỉ
chỉ chính
chỉ có
chỉ là
chỉ tên
chỉn
chị
chị bộ
chị ấy
chịu
chịu chưa
chịu lời
chịu tốt
chịu ăn
chọn
chọn bên
chọn ra
chốc chốc
chớ
chớ chi
chớ gì
chớ không
chớ kể
chớ như
chợt
chợt nghe
chợt nhìn
chủn
chứ
chứ ai
chứ còn
chứ gì
chứ không
chứ không phải
chứ lại
chứ lị
chứ như
chứ sao
coi bộ
coi mòi
con
con con
con dạ
con nhà
con tính
cu cậu
cuối
cuối cùng
cuối điểm
cuốn
cuộc
càng
càng càng
càng hay
cá nhân
các
các cậu
cách
cách bức
cách không
cách nhau
cách đều
cái
cái gì
cái họ
cái đã
cái đó
cái ấy
câu hỏi
cây
cây nước
còn
còn như
còn nữa
còn thời gian
còn về
có
có ai
có chuyện
có chăng
có chăng là
có chứ
có cơ
có dễ
có họ
có khi
có ngày
có người
có nhiều
có nhà
có phải
có số
có tháng
có thế
có thể
có vẻ
có ý
có ăn
có điều
có điều kiện
có đáng
có đâu
có được
cóc khô
cô
cô mình
cô quả
cô tăng
cô ấy
công nhiên
cùng
cùng chung
cùng cực
cùng nhau
cùng tuổi
cùng tột
cùng với
cùng ăn
căn
căn cái
căn cắt
căn tính
cũng
cũng như
cũng nên
cũng thế
cũng vậy
cũng vậy thôi
cũng được
cơ
cơ chỉ
cơ chừng
cơ cùng
cơ dẫn
cơ hồ
cơ hội
cơ mà
cơn
cả
cả nghe
cả nghĩ
cả ngày
cả người
cả nhà
cả năm
cả thảy
cả thể
cả tin
cả ăn
cả đến
cảm thấy
cảm ơn
cấp
cấp số
cấp trực tiếp
cần
cần cấp
cần gì
cần số
cật lực
cật sức
cậu
cổ lai
cụ thể
cụ thể là
cụ thể như
của
của ngọt
của tin
cứ
cứ như
cứ việc
cứ điểm
cực lực
do
do vì
do vậy
do đó
duy
duy chỉ
duy có
dài
dài lời
dài ra
dành
dành dành
dào
dì
dù
dù cho
dù dì
dù gì
dù rằng
dù sao
dùng
dùng cho
dùng hết
dùng làm
dùng đến
dưới
dưới nước
dạ
dạ bán
dạ con
dạ dài
dạ dạ
dạ khách
dần dà
dần dần
dầu sao
dẫn
dẫu
dẫu mà
dẫu rằng
dẫu sao
dễ
dễ dùng
dễ gì
dễ khiến
dễ nghe
dễ ngươi
dễ như chơi
dễ sợ
dễ sử dụng
dễ thường
dễ thấy
dễ ăn
dễ đâu
dở chừng
dữ
dữ cách
em
em em
giá trị
giá trị thực tế
giảm
giảm chính
giảm thấp
giảm thế
giống
giống người
giống nhau
giống như
giờ
giờ lâu
giờ này
giờ đi
giờ đây
giờ đến
giữ
giữ lấy
giữ ý
giữa
giữa lúc
gây
gây cho
gây giống
gây ra
gây thêm
gì
gì gì
gì đó
gần
gần bên
gần hết
gần ngày
gần như
gần xa
gần đây
gần đến
gặp
gặp khó khăn
gặp phải
gồm
hay
hay biết
hay hay
hay không
hay là
hay làm
hay nhỉ
hay nói
hay sao
hay tin
hay đâu
hiểu
hiện nay
hiện tại
hoàn toàn
hoặc
hoặc là
hãy
hãy còn
hơn
hơn cả
hơn hết
hơn là
hơn nữa
hơn trước
hầu hết
hết
hết chuyện
hết cả
hết của
hết nói
hết ráo
hết rồi
hết ý
họ
họ gần
họ xa
hỏi
hỏi lại
hỏi xem
hỏi xin
hỗ trợ
khi
khi khác
khi không
khi nào
khi nên
khi trước
khiến
khoảng
khoảng cách
khoảng không
khá
khá tốt
khác
khác gì
khác khác
khác nhau
khác nào
khác thường
khác xa
khách
khó
khó biết
khó chơi
khó khăn
khó làm
khó mở
khó nghe
khó nghĩ
khó nói
khó thấy
khó tránh
không
không ai
không bao giờ
không bao lâu
không biết
không bán
không chỉ
không còn
không có
không có gì
không cùng
không cần
không cứ
không dùng
không gì
không hay
không khỏi
không kể
không ngoài
không nhận
không những
không phải
không phải không
không thể
không tính
không điều kiện
không được
không đầy
không để
khẳng định
khỏi
khỏi nói
kể
kể cả
kể như
kể tới
kể từ
liên quan
loại
loại từ
luôn
luôn cả
luôn luôn
luôn tay
là
là cùng
là là
là nhiều
là phải
là thế nào
là vì
là ít
làm
làm bằng
làm cho
làm dần dần
làm gì
làm lòng
làm lại
làm lấy
làm mất
làm ngay
làm như
làm nên
làm ra
làm riêng
làm sao
làm theo
làm thế nào
làm tin
làm tôi
làm tăng
làm tại
làm tắp lự
làm vì
làm đúng
làm được
lâu
lâu các
lâu lâu
lâu nay
lâu ngày
lên
lên cao
lên cơn
lên mạnh
lên ngôi
lên nước
lên số
lên xuống
lên đến
lòng
lòng không
lúc
lúc khác
lúc lâu
lúc nào
lúc này
lúc sáng
lúc trước
lúc đi
lúc đó
lúc đến
lúc ấy
lý do
lượng
lượng cả
lượng số
lượng từ
lại
lại bộ
lại cái
lại còn
lại giống
lại làm
lại người
lại nói
lại nữa
lại quả
lại thôi
lại ăn
lại đây
lấy
lấy có
lấy cả
lấy giống
lấy làm
lấy lý do
lấy lại
lấy ra
lấy ráo
lấy sau
lấy số
lấy thêm
lấy thế
lấy vào
lấy xuống
lấy được
lấy để
lần
lần khác
lần lần
lần nào
lần này
lần sang
lần sau
lần theo
lần trước
lần tìm
lớn
lớn lên
lớn nhỏ
lời
lời chú
lời nói
mang
mang lại
mang mang
mang nặng
mang về
muốn
mà
mà cả
mà không
mà lại
mà thôi
mà vẫn
mình
mạnh
mất
mất còn
mọi
mọi giờ
mọi khi
mọi lúc
mọi người
mọi nơi
mọi sự
mọi thứ
mọi việc
mối
mỗi
mỗi lúc
mỗi lần
mỗi một
mỗi ngày
mỗi người
một
một cách
một cơn
một khi
một lúc
một số
một vài
một ít
mới
mới hay
mới rồi
mới đây
mở
mở mang
mở nước
mở ra
mợ
mức
nay
ngay
ngay bây giờ
ngay cả
ngay khi
ngay khi đến
ngay lúc
ngay lúc này
ngay lập tức
ngay thật
ngay tức khắc
ngay tức thì
ngay từ
nghe
nghe chừng
nghe hiểu
nghe không
nghe lại
nghe nhìn
nghe như
nghe nói
nghe ra
nghe rõ
nghe thấy
nghe tin
nghe trực tiếp
nghe đâu
nghe đâu như
nghe được
nghen
nghiễm nhiên
nghĩ
nghĩ lại
nghĩ ra
nghĩ tới
nghĩ xa
nghĩ đến
nghỉm
ngoài
ngoài này
ngoài ra
ngoài xa
ngoải
nguồn
ngày
ngày càng
ngày cấp
ngày giờ
ngày ngày
ngày nào
ngày này
ngày nọ
ngày qua
ngày rày
ngày tháng
ngày xưa
ngày xửa
ngày đến
ngày ấy
ngôi
ngôi nhà
ngôi thứ
ngõ hầu
ngăn ngắt
ngươi
người
người hỏi
người khác
người khách
người mình
người nghe
người người
người nhận
ngọn
ngọn nguồn
ngọt
ngồi
ngồi bệt
ngồi không
ngồi sau
ngồi trệt
ngộ nhỡ
nhanh
nhanh lên
nhanh tay
nhau
nhiên hậu
nhiều
nhiều ít
nhiệt liệt
nhung nhăng
nhà
nhà chung
nhà khó
nhà làm
nhà ngoài
nhà ngươi
nhà tôi
nhà việc
nhân dịp
nhân tiện
nhé
nhìn
nhìn chung
nhìn lại
nhìn nhận
nhìn theo
nhìn thấy
nhìn xuống
nhóm
nhón nhén
như
như ai
như chơi
như không
như là
như nhau
như quả
như sau
như thường
như thế
như thế nào
như thể
như trên
như trước
như tuồng
như vậy
như ý
nhưng
nhưng mà
nhược bằng
nhất
nhất loạt
nhất luật
nhất là
nhất mực
nhất nhất
nhất quyết
nhất sinh
nhất thiết
nhất thì
nhất tâm
nhất tề
nhất đán
nhất định
nhận
nhận biết
nhận họ
nhận làm
nhận nhau
nhận ra
nhận thấy
nhận việc
nhận được
nhằm
nhằm khi
nhằm lúc
nhằm vào
nhằm để
nhỉ
nhỏ
nhỏ người
nhớ
nhớ bập bõm
nhớ lại
nhớ lấy
nhớ ra
nhờ
nhờ chuyển
nhờ có
nhờ nhờ
nhờ đó
nhỡ ra
những
những ai
những khi
những là
những lúc
những muốn
những như
nào
nào cũng
nào hay
nào là
nào phải
nào đâu
nào đó
này
này nọ
nên
nên chi
nên chăng
nên làm
nên người
nên tránh
nó
nóc
nói
nói bông
nói chung
nói khó
nói là
nói lên
nói lại
nói nhỏ
nói phải
nói qua
nói ra
nói riêng
nói rõ
nói thêm
nói thật
nói toẹt
nói trước
nói tốt
nói với
nói xa
nói ý
nói đến
nói đủ
năm
năm tháng
nơi
nơi nơi
nước
nước bài
nước cùng
nước lên
nước nặng
nước quả
nước xuống
nước ăn
nước đến
nấy
nặng
nặng căn
nặng mình
nặng về
nếu
nếu có
nếu cần
nếu không
nếu mà
nếu như
nếu thế
nếu vậy
nếu được
nền
nọ
nớ
nức nở
nữa
nữa khi
nữa là
nữa rồi
oai oái
oái
pho
phè
phè phè
phía
phía bên
phía bạn
phía dưới
phía sau
phía trong
phía trên
phía trước
phóc
phót
phù hợp
phăn phắt
phương chi
phải
phải biết
phải chi
phải chăng
phải cách
phải cái
phải giờ
phải khi
phải không
phải lại
phải lời
phải người
phải như
phải rồi
phải tay
phần
phần lớn
phần nhiều
phần nào
phần sau
phần việc
phắt
phỉ phui
phỏng
phỏng như
phỏng nước
phỏng theo
phỏng tính
phốc
phụt
phứt
qua
qua chuyện
qua khỏi
qua lại
qua lần
qua ngày
qua tay
qua thì
qua đi
quan trọng
quan trọng vấn đề
quan tâm
quay
quay bước
quay lại
quay số
quay đi
quá
quá bán
quá bộ
quá giờ
quá lời
quá mức
quá nhiều
quá tay
quá thì
quá tin
quá trình
quá tuổi
quá đáng
quá ư
quả
quả là
quả thật
quả thế
quả vậy
quận
ra
ra bài
ra bộ
ra chơi
ra gì
ra lại
ra lời
ra ngôi
ra người
ra sao
ra tay
ra vào
ra ý
ra điều
ra đây
ren rén
riu ríu
riêng
riêng từng
riệt
rày
ráo
ráo cả
ráo nước
ráo trọi
rén
rén bước
rích
rón rén
rõ
rõ là
rõ thật
rút cục
răng
răng răng
rất
rất lâu
rằng
rằng là
rốt cuộc
rốt cục
rồi
rồi nữa
rồi ra
rồi sao
rồi sau
rồi tay
rồi thì
rồi xem
rồi đây
rứa
sa sả
sang
sang năm
sang sáng
sang tay
sao
sao bản
sao bằng
sao cho
sao vậy
sao đang
sau
sau chót
sau cuối
sau cùng
sau hết
sau này
sau nữa
sau sau
sau đây
sau đó
so
so với
song le
suýt
suýt nữa
sáng
sáng ngày
sáng rõ
sáng thế
sáng ý
sì
sì sì
sất
sắp
sắp đặt
sẽ
sẽ biết
sẽ hay
số
số cho biết
số cụ thể
số loại
số là
số người
số phần
số thiếu
sốt sột
sớm
sớm ngày
sở dĩ
sử dụng
sự
sự thế
sự việc
tanh
tanh tanh
tay
tay quay
tha hồ
tha hồ chơi
tha hồ ăn
than ôi
thanh
thanh ba
thanh chuyển
thanh không
thanh thanh
thanh tính
thanh điều kiện
thanh điểm
thay đổi
thay đổi tình trạng
theo
theo bước
theo như
theo tin
thi thoảng
thiếu
thiếu gì
thiếu điểm
thoạt
thoạt nghe
thoạt nhiên
thoắt
thuần
thuần ái
thuộc
thuộc bài
thuộc cách
thuộc lại
thuộc từ
thà
thà là
thà rằng
thành ra
thành thử
thái quá
tháng
tháng ngày
tháng năm
tháng tháng
thêm
thêm chuyện
thêm giờ
thêm vào
thì
thì giờ
thì là
thì phải
thì ra
thì thôi
thình lình
thích
thích cứ
thích thuộc
thích tự
thích ý
thím
thôi
thôi việc
thúng thắng
thương ôi
thường
thường bị
thường hay
thường khi
thường số
thường sự
thường thôi
thường thường
thường tính
thường tại
thường xuất hiện
thường đến
thảo hèn
thảo nào
thấp
thấp cơ
thấp thỏm
thấp xuống
thấy
thấy tháng
thẩy
thậm
thậm chí
thậm cấp
thậm từ
thật
thật chắc
thật là
thật lực
thật quả
thật ra
thật sự
thật thà
thật tốt
thật vậy
thế
thế chuẩn bị
thế là
thế lại
thế mà
thế nào
thế nên
thế ra
thế sự
thế thì
thế thôi
thế thường
thế thế
thế à
thế đó
thếch
thỉnh thoảng
thỏm
thốc
thốc tháo
thốt
thốt nhiên
thốt nói
thốt thôi
thộc
thời gian
thời gian sử dụng
thời gian tính
thời điểm
thục mạng
thứ
thứ bản
thứ đến
thửa
thực hiện
thực hiện đúng
thực ra
thực sự
thực tế
thực vậy
tin
tin thêm
tin vào
tiếp theo
tiếp tục
tiếp đó
tiện thể
toà
toé khói
toẹt
trong
trong khi
trong lúc
trong mình
trong ngoài
trong này
trong số
trong vùng
trong đó
trong ấy
tránh
tránh khỏi
tránh ra
tránh tình trạng
tránh xa
trên
trên bộ
trên dưới
trước
trước hết
trước khi
trước kia
trước nay
trước ngày
trước nhất
trước sau
trước tiên
trước tuổi
trước đây
trước đó
trả
trả của
trả lại
trả ngay
trả trước
trếu tráo
trển
trệt
trệu trạo
trỏng
trời đất ơi
trở thành
trừ phi
trực tiếp
trực tiếp làm
tuy
tuy có
tuy là
tuy nhiên
tuy rằng
tuy thế
tuy vậy
tuy đã
tuyệt nhiên
tuần tự
tuốt luốt
tuốt tuồn tuột
tuốt tuột
tuổi
tuổi cả
tuổi tôi
tà tà
tên
tên chính
tên cái
tên họ
tên tự
tênh
tênh tênh
tìm
tìm bạn
tìm cách
tìm hiểu
tìm ra
tìm việc
tình trạng
tính
tính cách
tính căn
tính người
tính phỏng
tính từ
tít mù
tò te
tôi
tôi con
tông tốc
tù tì
tăm tắp
tăng
tăng chúng
tăng cấp
tăng giảm
tăng thêm
tăng thế
tại
tại lòng
tại nơi
tại sao
tại tôi
tại vì
tại đâu
tại đây
tại đó
tạo
tạo cơ hội
tạo nên
tạo ra
tạo ý
tạo điều kiện
tấm
tấm bản
tấm các
tấn
tấn tới
tất cả
tất cả bao nhiêu
tất thảy
tất tần tật
tất tật
tập trung
tắp
tắp lự
tắp tắp
tọt
tỏ ra
tỏ vẻ
tốc tả
tối ư
tốt
tốt bạn
tốt bộ
tốt hơn
tốt mối
tốt ngày
tột
tột cùng
tớ
tới
tới gần
tới mức
tới nơi
tới thì
tức thì
tức tốc
từ
từ căn
từ giờ
từ khi
từ loại
từ nay
từ thế
từ tính
từ tại
từ từ
từ ái
từ điều
từ đó
từ ấy
từng
từng cái
từng giờ
từng nhà
từng phần
từng thời gian
từng đơn vị
từng ấy
tự
tự cao
tự khi
tự lượng
tự tính
tự tạo
tự vì
tự ý
tự ăn
tựu trung
veo
veo veo
việc
việc gì
vung thiên địa
vung tàn tán
vung tán tàn
và
vài
vài ba
vài người
vài nhà
vài nơi
vài tên
vài điều
vào
vào gặp
vào khoảng
vào lúc
vào vùng
vào đến
vâng
vâng chịu
vâng dạ
vâng vâng
vâng ý
vèo
vèo vèo
vì
vì chưng
vì rằng
vì sao
vì thế
vì vậy
ví bằng
ví dù
ví phỏng
ví thử
vô hình trung
vô kể
vô luận
vô vàn
vùng
vùng lên
vùng nước
văng tê
vượt
vượt khỏi
vượt quá
vạn nhất
vả chăng
vả lại
vấn đề
vấn đề quan trọng
vẫn
vẫn thế
vậy
vậy là
vậy mà
vậy nên
vậy ra
vậy thì
vậy ư
về
về không
về nước
về phần
về sau
về tay
vị trí
vị tất
vốn dĩ
với
với lại
với nhau
vở
vụt
vừa
vừa khi
vừa lúc
vừa mới
vừa qua
vừa rồi
vừa vừa
xa
xa cách
xa gần
xa nhà
xa tanh
xa tắp
xa xa
xa xả
xem
xem lại
xem ra
xem số
xin
xin gặp
xin vâng
xiết bao
xon xón
xoành xoạch
xoét
xoẳn
xoẹt
xuất hiện
xuất kì bất ý
xuất kỳ bất ý
xuể
xuống
xăm xúi
xăm xăm
xăm xắm
xảy ra
xềnh xệch
xệp
xử lý
yêu cầu
à
à này
à ơi
ào
ào vào
ào ào
á
á à
ái
ái chà
ái dà
áng
áng như
âu là
ít
ít biết
ít có
ít hơn
ít khi
ít lâu
ít nhiều
ít nhất
ít nữa
ít quá
ít ra
ít thôi
ít thấy
ô hay
ô hô
ô kê
ô kìa
ôi chao
ôi thôi
ông
ông nhỏ
ông tạo
ông từ
ông ấy
ông ổng
úi
úi chà
úi dào
ý
ý chừng
ý da
ý hoặc
ăn
ăn chung
ăn chắc
ăn chịu
ăn cuộc
ăn hết
ăn hỏi
ăn làm
ăn người
ăn ngồi
ăn quá
ăn riêng
ăn sáng
ăn tay
ăn trên
ăn về
đang
đang tay
đang thì
điều
điều gì
điều kiện
điểm
điểm chính
điểm gặp
điểm đầu tiên
đành đạch
đáng
đáng kể
đáng lí
đáng lý
đáng lẽ
đáng số
đánh giá
đánh đùng
đáo để
đâu
đâu có
đâu cũng
đâu như
đâu nào
đâu phải
đâu đâu
đâu đây
đâu đó
đây
đây này
đây rồi
đây đó
đã
đã hay
đã không
đã là
đã lâu
đã thế
đã vậy
đã đủ
đó
đó đây
đúng
đúng ngày
đúng ra
đúng tuổi
đúng với
đơn vị
đưa
đưa cho
đưa chuyện
đưa em
đưa ra
đưa tay
đưa tin
đưa tới
đưa vào
đưa về
đưa xuống
đưa đến
được
được cái
được lời
được nước
được tin
đại loại
đại nhân
đại phàm
đại để
đạt
đảm bảo
đầu tiên
đầy
đầy năm
đầy phè
đầy tuổi
đặc biệt
đặt
đặt làm
đặt mình
đặt mức
đặt ra
đặt trước
đặt để
đến
đến bao giờ
đến cùng
đến cùng cực
đến cả
đến giờ
đến gần
đến hay
đến khi
đến lúc
đến lời
đến nay
đến ngày
đến nơi
đến nỗi
đến thì
đến thế
đến tuổi
đến xem
đến điều
đến đâu
đều
đều bước
đều nhau
đều đều
để
để cho
để giống
để không
để lòng
để lại
để mà
để phần
để được
để đến nỗi
đối với
đồng thời
đủ
đủ dùng
đủ nơi
đủ số
đủ điều
đủ điểm
ơ
ơ hay
ơ kìa
ơi
ơi là
ư
ạ
ạ ơi
ấy
ấy là
ầu ơ
ắt
ắt hẳn
ắt là
ắt phải
ắt thật
ối dào
ối giời
ối giời ơi
ồ
ồ ồ
ổng
ớ
ớ này
ờ
ờ ờ
ở
ở lại
ở như
ở nhờ
ở năm
ở trên
ở vào
ở đây
ở đó
ở được
ủa
ứ hự
ứ ừ
ừ
ừ nhé
ừ thì
ừ ào
ừ ừ
ử
`,
}
