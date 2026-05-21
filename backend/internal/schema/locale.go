package schema

import (
	"fmt"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
)

// LocaleInfo holds locale-specific data for realistic fake data generation.
type LocaleInfo struct {
	FirstNames  []string
	LastNames   []string
	Cities      []string
	States      []string
	Country     string
	PhonePrefix string
	PhoneDigits int // digits after prefix
	ZipFormat   string
}

// localeDB maps BCP-47 locale tags to locale data.
// Locale tags are lower-cased on lookup so both "en-IN" and "en-in" work.
var localeDB = map[string]*LocaleInfo{
	"en-in": {
		FirstNames: []string{
			"Aarav", "Aditi", "Akash", "Anika", "Ananya", "Arjun", "Aryan",
			"Divya", "Ishaan", "Kavya", "Kritika", "Meera", "Mohan", "Neha",
			"Nikhil", "Pooja", "Priya", "Rahul", "Rajesh", "Riya", "Rohit",
			"Sandeep", "Shruti", "Sneha", "Suresh", "Tanvi", "Vikram", "Vikas",
			"Vijay", "Vivek", "Aditya", "Anjali", "Deepak", "Gaurav", "Harsha",
			"Isha", "Kiran", "Lakshmi", "Manish", "Nandita", "Omkar", "Pankaj",
			"Rekha", "Sameer", "Smita", "Tarun", "Uma", "Varun", "Yogesh", "Zara",
		},
		LastNames: []string{
			"Sharma", "Singh", "Verma", "Patel", "Kumar", "Gupta", "Joshi",
			"Rao", "Nair", "Mehta", "Shah", "Reddy", "Mishra", "Agarwal",
			"Pandey", "Tiwari", "Malhotra", "Chaudhary", "Kapoor", "Bose",
			"Chatterjee", "Mukherjee", "Pillai", "Iyer", "Menon", "Desai",
			"Saxena", "Choudhury", "Bhatia", "Naik",
		},
		Cities: []string{
			"Mumbai", "Delhi", "Bangalore", "Chennai", "Kolkata", "Hyderabad",
			"Pune", "Ahmedabad", "Jaipur", "Lucknow", "Surat", "Kanpur",
			"Nagpur", "Indore", "Bhopal", "Visakhapatnam", "Patna", "Vadodara",
			"Ludhiana", "Agra", "Nashik", "Coimbatore", "Madurai", "Kochi",
			"Chandigarh", "Mysore", "Bhubaneswar", "Thiruvananthapuram", "Guwahati",
		},
		States: []string{
			"Maharashtra", "Karnataka", "Tamil Nadu", "Delhi", "West Bengal",
			"Telangana", "Gujarat", "Rajasthan", "Uttar Pradesh", "Madhya Pradesh",
			"Punjab", "Bihar", "Andhra Pradesh", "Kerala", "Odisha",
			"Assam", "Uttarakhand", "Himachal Pradesh", "Goa", "Jharkhand",
		},
		Country:     "India",
		PhonePrefix: "+91",
		PhoneDigits: 10,
		ZipFormat:   "######",
	},
	"ja-jp": {
		FirstNames: []string{
			"Hiroshi", "Yuki", "Kenji", "Akiko", "Takashi", "Yuko", "Daisuke",
			"Naomi", "Haru", "Sakura", "Ren", "Aoi", "Sota", "Mia", "Kaito",
			"Nana", "Riku", "Haruka", "Kenta", "Misaki", "Yuta", "Saki",
			"Shota", "Rina", "Daiki", "Yui", "Ryota", "Miyu", "Taro", "Hanako",
		},
		LastNames: []string{
			"Sato", "Suzuki", "Tanaka", "Watanabe", "Ito", "Yamamoto",
			"Nakamura", "Kobayashi", "Kato", "Yoshida", "Yamada", "Sasaki",
			"Yamaguchi", "Matsumoto", "Inoue", "Kimura", "Hayashi", "Shimizu",
			"Yamazaki", "Mori",
		},
		Cities: []string{
			"Tokyo", "Osaka", "Nagoya", "Sapporo", "Fukuoka", "Kobe",
			"Kyoto", "Kawasaki", "Saitama", "Hiroshima", "Sendai", "Yokohama",
			"Chiba", "Kitakyushu", "Sakai", "Niigata", "Hamamatsu", "Kumamoto",
		},
		States: []string{
			"Tokyo", "Osaka", "Kanagawa", "Aichi", "Saitama", "Chiba",
			"Hyogo", "Hokkaido", "Fukuoka", "Shizuoka", "Ibaraki", "Hiroshima",
		},
		Country:     "Japan",
		PhonePrefix: "+81",
		PhoneDigits: 10,
		ZipFormat:   "###-####",
	},
	"de-de": {
		FirstNames: []string{
			"Hans", "Wolfgang", "Klaus", "Dieter", "Michael", "Thomas", "Andreas",
			"Stefan", "Petra", "Sabine", "Monika", "Claudia", "Anna", "Laura",
			"Lena", "Emma", "Maximilian", "Alexander", "Julian", "Lukas",
			"Felix", "Sophie", "Marie", "Hannah", "Leonie",
		},
		LastNames: []string{
			"Müller", "Schmidt", "Schneider", "Fischer", "Weber", "Meyer",
			"Wagner", "Becker", "Schulz", "Hoffmann", "Schäfer", "Koch",
			"Bauer", "Richter", "Klein", "Wolf", "Schröder", "Neumann",
			"Schwarz", "Zimmermann",
		},
		Cities: []string{
			"Berlin", "Hamburg", "Munich", "Cologne", "Frankfurt", "Stuttgart",
			"Düsseldorf", "Dortmund", "Essen", "Leipzig", "Bremen", "Dresden",
			"Hanover", "Nuremberg", "Duisburg", "Bochum", "Wuppertal", "Bonn",
		},
		States: []string{
			"Bavaria", "North Rhine-Westphalia", "Baden-Württemberg", "Berlin",
			"Hamburg", "Hesse", "Saxony", "Lower Saxony", "Rhineland-Palatinate",
			"Brandenburg", "Thuringia", "Saxony-Anhalt",
		},
		Country:     "Germany",
		PhonePrefix: "+49",
		PhoneDigits: 10,
		ZipFormat:   "#####",
	},
	"fr-fr": {
		FirstNames: []string{
			"Pierre", "Jean", "Louis", "Michel", "René", "André", "François",
			"Philippe", "Jacques", "Antoine", "Marie", "Sophie", "Julie",
			"Laura", "Emma", "Camille", "Léa", "Manon", "Inès", "Lucie",
			"Lucas", "Hugo", "Nathan", "Théo", "Mathis",
		},
		LastNames: []string{
			"Martin", "Bernard", "Thomas", "Petit", "Robert", "Richard",
			"Durand", "Dubois", "Moreau", "Laurent", "Simon", "Michel",
			"Lefebvre", "Leroy", "Roux", "David", "Bertrand", "Morel",
			"Fournier", "Girard",
		},
		Cities: []string{
			"Paris", "Marseille", "Lyon", "Toulouse", "Nice", "Nantes",
			"Strasbourg", "Montpellier", "Bordeaux", "Lille", "Rennes",
			"Reims", "Le Havre", "Saint-Étienne", "Toulon", "Grenoble",
		},
		States: []string{
			"Île-de-France", "Provence-Alpes-Côte d'Azur", "Auvergne-Rhône-Alpes",
			"Occitanie", "Nouvelle-Aquitaine", "Grand Est", "Hauts-de-France",
			"Brittany", "Normandy", "Pays de la Loire",
		},
		Country:     "France",
		PhonePrefix: "+33",
		PhoneDigits: 9,
		ZipFormat:   "#####",
	},
	"zh-cn": {
		FirstNames: []string{
			"Wei", "Fang", "Xia", "Yu", "Jing", "Li", "Mei", "Ming", "Hao",
			"Jun", "Ting", "Xin", "Lei", "Bo", "Yan", "Hui", "Tao", "Qian",
			"Rui", "Shu", "Xiao", "Yun", "Zhen", "Hua", "Peng",
		},
		LastNames: []string{
			"Wang", "Li", "Zhang", "Liu", "Chen", "Yang", "Zhao", "Wu",
			"Zhou", "Sun", "Xu", "Ma", "Zhu", "Hu", "Guo", "He", "Lin",
			"Gao", "Luo", "Zheng",
		},
		Cities: []string{
			"Shanghai", "Beijing", "Guangzhou", "Shenzhen", "Chengdu",
			"Tianjin", "Wuhan", "Chongqing", "Xi'an", "Hangzhou",
			"Nanjing", "Harbin", "Shenyang", "Qingdao", "Zhengzhou",
		},
		States: []string{
			"Beijing", "Shanghai", "Guangdong", "Sichuan", "Zhejiang",
			"Jiangsu", "Shandong", "Henan", "Hubei", "Hunan",
			"Liaoning", "Shaanxi", "Fujian",
		},
		Country:     "China",
		PhonePrefix: "+86",
		PhoneDigits: 11,
		ZipFormat:   "######",
	},
	"pt-br": {
		FirstNames: []string{
			"João", "Maria", "José", "Ana", "Carlos", "Paulo", "Pedro",
			"Luiz", "Marcos", "Lucas", "Gabriel", "Rafael", "Matheus",
			"Guilherme", "Fernanda", "Juliana", "Amanda", "Leticia", "Larissa", "Isabela",
		},
		LastNames: []string{
			"Silva", "Santos", "Oliveira", "Souza", "Rodrigues", "Ferreira",
			"Alves", "Pereira", "Lima", "Gomes", "Costa", "Ribeiro",
			"Martins", "Carvalho", "Almeida", "Lopes", "Sousa", "Fernandes",
			"Vieira", "Barbosa",
		},
		Cities: []string{
			"São Paulo", "Rio de Janeiro", "Brasília", "Salvador", "Fortaleza",
			"Belo Horizonte", "Manaus", "Curitiba", "Recife", "Porto Alegre",
			"Goiânia", "Belém", "São Luís", "Maceió", "Natal",
		},
		States: []string{
			"São Paulo", "Rio de Janeiro", "Minas Gerais", "Bahia", "Paraná",
			"Rio Grande do Sul", "Pernambuco", "Ceará", "Pará", "Goiás",
			"Amazonas", "Santa Catarina", "Maranhão", "Mato Grosso", "Espírito Santo",
		},
		Country:     "Brazil",
		PhonePrefix: "+55",
		PhoneDigits: 11,
		ZipFormat:   "#####-###",
	},
	"es-es": {
		FirstNames: []string{
			"Alejandro", "David", "Daniel", "Carlos", "Pablo", "Javier",
			"José", "Antonio", "Manuel", "Pedro", "María", "Laura", "Ana",
			"Sara", "Marta", "Lucía", "Carmen", "Rosa", "Isabel", "Elena",
		},
		LastNames: []string{
			"García", "Rodríguez", "González", "Fernández", "López",
			"Martínez", "Sánchez", "Pérez", "Gómez", "Díaz", "Torres",
			"Ramírez", "Ruiz", "Flores", "Jiménez", "Hernández", "Moreno",
			"Muñoz", "Álvarez", "Romero",
		},
		Cities: []string{
			"Madrid", "Barcelona", "Valencia", "Seville", "Zaragoza",
			"Málaga", "Murcia", "Palma", "Las Palmas", "Bilbao",
			"Alicante", "Córdoba", "Valladolid", "Vigo", "Gijón",
		},
		States: []string{
			"Andalusia", "Catalonia", "Community of Madrid", "Valencian Community",
			"Galicia", "Castile and León", "Basque Country", "Canary Islands",
			"Castile-La Mancha", "Aragon", "Extremadura", "Asturias",
		},
		Country:     "Spain",
		PhonePrefix: "+34",
		PhoneDigits: 9,
		ZipFormat:   "#####",
	},
	"ko-kr": {
		FirstNames: []string{
			"Minjun", "Jiyeon", "Sohee", "Jaeho", "Seulki", "Hyunwoo",
			"Eunjin", "Youngsu", "Jinseo", "Nayeon", "Taehyun", "Minji",
			"Seojun", "Yuna", "Junho", "Sooji", "Hyunjin", "Jiyoo",
			"Minhee", "Donghyun",
		},
		LastNames: []string{
			"Kim", "Lee", "Park", "Choi", "Jung", "Kang", "Cho",
			"Yoon", "Jang", "Lim", "Han", "Oh", "Seo", "Shin", "Kwon",
			"Hwang", "Ahn", "Song", "Baek", "Ryu",
		},
		Cities: []string{
			"Seoul", "Busan", "Incheon", "Daegu", "Daejeon", "Gwangju",
			"Suwon", "Ulsan", "Changwon", "Goyang", "Yongin", "Seongnam",
		},
		States: []string{
			"Seoul", "Busan", "Incheon", "Daegu", "Daejeon", "Gwangju",
			"Gyeonggi", "Gyeongnam", "Gyeongbuk", "Jeonnam", "Jeonbuk",
			"Chungnam", "Chungbuk", "Gangwon", "Jeju",
		},
		Country:     "South Korea",
		PhonePrefix: "+82",
		PhoneDigits: 10,
		ZipFormat:   "#####",
	},
	"ar-sa": {
		FirstNames: []string{
			"Mohammed", "Ahmed", "Ali", "Omar", "Abdullah", "Ibrahim",
			"Hassan", "Hussain", "Khalid", "Faisal", "Fatima", "Aisha",
			"Sara", "Noor", "Layla", "Maryam", "Reem", "Hessa", "Noura", "Rania",
		},
		LastNames: []string{
			"Al-Saud", "Al-Rashid", "Al-Otaibi", "Al-Shammari", "Al-Zahrani",
			"Al-Qahtani", "Al-Maliki", "Al-Harbi", "Al-Ghamdi", "Al-Dosari",
			"Al-Mutairi", "Al-Aziz", "Al-Hamdan", "Al-Jaber", "Al-Khalid",
		},
		Cities: []string{
			"Riyadh", "Jeddah", "Mecca", "Medina", "Dammam", "Taif",
			"Tabuk", "Buraidah", "Khobar", "Abha", "Hail", "Najran",
		},
		States: []string{
			"Riyadh", "Mecca", "Eastern Province", "Medina", "Asir",
			"Jazan", "Najran", "Tabuk", "Hail", "Northern Borders",
		},
		Country:     "Saudi Arabia",
		PhonePrefix: "+966",
		PhoneDigits: 9,
		ZipFormat:   "#####",
	},
	"en-gb": {
		FirstNames: []string{
			"Oliver", "George", "Harry", "Jack", "Noah", "Charlie", "Jacob",
			"Alfie", "Freddie", "Oscar", "Amelia", "Olivia", "Isla", "Emily",
			"Ava", "Poppy", "Isabella", "Sophie", "Mia", "Freya",
		},
		LastNames: []string{
			"Smith", "Jones", "Williams", "Taylor", "Brown", "Davies", "Evans",
			"Wilson", "Thomas", "Roberts", "Johnson", "Lewis", "Walker",
			"Robinson", "Wood", "Thompson", "White", "Watson", "Jackson", "Wright",
		},
		Cities: []string{
			"London", "Birmingham", "Leeds", "Glasgow", "Sheffield", "Bradford",
			"Manchester", "Edinburgh", "Liverpool", "Bristol", "Cardiff",
			"Leicester", "Coventry", "Nottingham", "Newcastle", "Belfast",
		},
		States: []string{
			"England", "Scotland", "Wales", "Northern Ireland",
			"Greater London", "West Midlands", "Greater Manchester",
			"West Yorkshire", "Lancashire", "Kent",
		},
		Country:     "United Kingdom",
		PhonePrefix: "+44",
		PhoneDigits: 10,
		ZipFormat:   "XX# #XX",
	},
}

// GetLocaleInfo returns locale data for the given BCP-47 tag, or nil if unsupported.
func GetLocaleInfo(locale string) *LocaleInfo {
	return localeDB[strings.ToLower(locale)]
}

// SupportedLocales returns all supported locale tags.
func SupportedLocales() []string {
	out := make([]string, 0, len(localeDB))
	for k := range localeDB {
		out = append(out, k)
	}
	return out
}

// generateLocalePhone builds a phone number for the given locale.
func generateLocalePhone(info *LocaleInfo, faker *gofakeit.Faker) string {
	switch info.ZipFormat {
	case "###-####": // ja-jp style
		return fmt.Sprintf("%s-%02d-%04d-%04d",
			info.PhonePrefix,
			faker.IntRange(10, 99),
			faker.IntRange(1000, 9999),
			faker.IntRange(1000, 9999),
		)
	case "XX# #XX": // en-gb (phone is separate)
		digits := buildDigitString(faker, info.PhoneDigits)
		return fmt.Sprintf("%s%s", info.PhonePrefix, digits)
	default:
		digits := buildDigitString(faker, info.PhoneDigits)
		return fmt.Sprintf("%s%s", info.PhonePrefix, digits)
	}
}

// generateLocaleZip builds a postal code matching the locale's format.
func generateLocaleZip(info *LocaleInfo, faker *gofakeit.Faker) string {
	switch info.ZipFormat {
	case "###-####": // Japan
		return fmt.Sprintf("%03d-%04d", faker.IntRange(100, 999), faker.IntRange(1000, 9999))
	case "#####-###": // Brazil
		return fmt.Sprintf("%05d-%03d", faker.IntRange(10000, 99999), faker.IntRange(100, 999))
	case "XX# #XX": // UK – simplified
		letters := "ABCDEFGHJKLMNPQRSTUVWXY"
		l := string(letters[faker.IntRange(0, len(letters)-1)])
		return fmt.Sprintf("%s%d%d %d%s%s",
			l,
			faker.IntRange(1, 9),
			faker.IntRange(1, 9),
			faker.IntRange(1, 9),
			string(letters[faker.IntRange(0, len(letters)-1)]),
			string(letters[faker.IntRange(0, len(letters)-1)]),
		)
	default:
		digits := info.ZipFormat
		if strings.Count(digits, "#") == 5 {
			return fmt.Sprintf("%05d", faker.IntRange(10000, 99999))
		}
		if strings.Count(digits, "#") == 6 {
			return fmt.Sprintf("%06d", faker.IntRange(100000, 999999))
		}
		return fmt.Sprintf("%d", faker.IntRange(10000, 999999))
	}
}

// buildDigitString generates a string of n random digits.
func buildDigitString(faker *gofakeit.Faker, n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteByte(byte('0' + faker.IntRange(0, 9)))
	}
	return b.String()
}
