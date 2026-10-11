# Android / Kotlin client template

Retrofit + OkHttp + kotlinx.serialization + coroutines. Adapt to Moshi/Gson if the project already uses one.

```kotlin
// build.gradle.kts (app)
implementation("com.squareup.retrofit2:retrofit:2.11.0")
implementation("com.squareup.okhttp3:okhttp:4.12.0")
implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:1.6.3")
implementation("com.jakewharton.retrofit:retrofit2-kotlinx-serialization-converter:1.0.0")
// plugins { kotlin("plugin.serialization") }

android {
    buildFeatures { buildConfig = true }
    defaultConfig {
        buildConfigField("String", "API_BASE_URL", "\"http://10.0.2.2:8080/api/v1/\"") // trailing slash required by Retrofit
        buildConfigField("String", "API_KEY", "\"\"")
    }
}
```

## Models

```kotlin
/** Path code of each game. Rules come from GET /games. */
enum class GameCode(val path: String) { FTN("ftn"), BL("bl"), PW("pw") }

@Serializable
data class Envelope<T>(val code: Int, val message: String, val data: T? = null)

@Serializable data class SpecialSpec(val name: String, val maxNumber: Int, val samePool: Boolean, val userPick: Boolean)
@Serializable data class Prize(val main: Int, val special: Boolean, val name: String)

@Serializable
data class GameInfo(
    val code: String, val name: String, val maxNumber: Int, val pickCount: Int, val bigFrom: Int,
    val special: SpecialSpec? = null, val prizes: List<Prize>,
)

@Serializable
data class Features(
    val sum: Int, val odd: Int, val even: Int, val big: Int, val small: Int, val prime: Int,
    val tenGroup: List<Int>,          // length varies by game
    val fullTenGroup: Boolean,
    val consecutive: String,          // open set: none, 2, 3, 22, 32, 222, ...
    val ac: Int, val span: Int, val mod3: List<Int>, val tailSum: Int,
    val sameTailGroups: Int, val maxGap: Int, val minGap: Int,
)

@Serializable
data class Draw(
    val period: String, val date: String, val yearIndex: Int,
    val numbers: List<Int>,
    val special: Int? = null,         // bl 特別號 / pw 第二區; null for ftn
    val features: Features,
) { val localDate: LocalDate get() = LocalDate.parse(date) }

@Serializable
data class DrawPage(val items: List<Draw>, val total: Int, val limit: Int, val offset: Int) {
    val hasMore get() = offset + items.size < total
}

@Serializable
data class NumberStat(val number: Int, val windowCount: Int, val totalCount: Int, val omission: Int, val maxOmission: Int)

@Serializable
data class NumberStats(
    val window: Int, val totalDraws: Int, val latestPeriod: String,
    val expectedWindow: Double, val numbers: List<NumberStat>,
    val specialExpectedWindow: Double? = null, val special: List<NumberStat>? = null,
)

@Serializable data class AnalyzeRequest(val numbers: List<Int>, val special: Int? = null)
@Serializable data class MatchCount(val main: Int, val special: Boolean, val count: Int, val prize: String? = null)

@Serializable
data class AnalyzeResult(
    val numbers: List<Int>, val special: Int? = null, val features: Features,
    val history: List<MatchCount>, val totalDraws: Int, val lastMatch3: Draw? = null,
)

@Serializable
data class GameHealth(val draws: Int, val latestPeriod: String? = null, val latestDate: String? = null, val loadedAt: String? = null)
@Serializable data class Health(val status: String, val games: Map<String, GameHealth>)
```

## Service and client

```kotlin
interface LotteryService {
    @GET("health") suspend fun health(): Response<Envelope<Health>>
    @GET("games") suspend fun games(): Response<Envelope<List<GameInfo>>>
    @GET("{game}/draws") suspend fun draws(
        @Path("game") game: String,
        @Query("limit") limit: Int = 20, @Query("offset") offset: Int = 0,
        @Query("from") from: String? = null, @Query("to") to: String? = null,
    ): Response<Envelope<DrawPage>>
    @GET("{game}/draws/latest") suspend fun latest(@Path("game") game: String): Response<Envelope<Draw>>
    @GET("{game}/draws/{period}") suspend fun draw(@Path("game") game: String, @Path("period") period: String): Response<Envelope<Draw>>
    @GET("{game}/numbers/stats") suspend fun numberStats(@Path("game") game: String, @Query("window") window: Int = 100): Response<Envelope<NumberStats>>
    @POST("{game}/analyze") suspend fun analyze(@Path("game") game: String, @Body body: AnalyzeRequest): Response<Envelope<AnalyzeResult>>
}

class ApiException(val code: Int, override val message: String, val httpStatus: Int? = null) : Exception(message) {
    val retryable get() = code == -1 || code == 50301
}

object LotteryApi {
    // explicitNulls = false: omit "special": null for ftn / bl requests
    private val json = Json { ignoreUnknownKeys = true; explicitNulls = false }

    private val client = OkHttpClient.Builder()
        .connectTimeout(10, TimeUnit.SECONDS)
        .readTimeout(15, TimeUnit.SECONDS)
        .addInterceptor { chain ->
            val req = chain.request().newBuilder().apply {
                if (BuildConfig.API_KEY.isNotEmpty()) header("X-API-Key", BuildConfig.API_KEY)
            }.build()
            chain.proceed(req)
        }
        .build()

    val service: LotteryService = Retrofit.Builder()
        .baseUrl(BuildConfig.API_BASE_URL)
        .client(client)
        .addConverterFactory(json.asConverterFactory("application/json".toMediaType()))
        .build()
        .create(LotteryService::class.java)

    /** Unwraps the envelope; error bodies are parsed too so the server's code/message survive. */
    suspend fun <T> call(block: suspend LotteryService.() -> Response<Envelope<T>>): T {
        val res = try { service.block() } catch (e: IOException) {
            throw ApiException(-1, e.message ?: "network error")
        }
        val env = res.body() ?: res.errorBody()?.string()?.let {
            runCatching { json.decodeFromString<Envelope<JsonElement>>(it) }.getOrNull()
        }?.let { throw ApiException(it.code, it.message, res.code()) }
            ?: throw ApiException(-1, "empty response", res.code())
        if (env.code != 0 || env.data == null) throw ApiException(env.code, env.message, res.code())
        return env.data
    }
}

fun validatePick(rules: GameInfo, numbers: List<Int>, special: Int?) {
    require(numbers.size == rules.pickCount && numbers.toSet().size == rules.pickCount &&
            numbers.all { it in 1..rules.maxNumber }) {
        "需選 ${rules.pickCount} 個不重複的 1~${rules.maxNumber} 號碼"
    }
    val sp = rules.special
    if (sp != null && sp.userPick) {
        require(special != null && special in 1..sp.maxNumber) { "需選 1 個 1~${sp.maxNumber} 的${sp.name}號碼" }
    } else {
        require(special == null) { "${rules.name}不需另外選特別號" }
    }
}

// usage in a repository / ViewModel
suspend fun loadLatest(game: GameCode): Draw = LotteryApi.call { latest(game.path) }
suspend fun analyze(rules: GameInfo, nums: List<Int>, special: Int? = null): AnalyzeResult {
    validatePick(rules, nums, special)
    return LotteryApi.call { analyze(rules.code, AnalyzeRequest(nums, special)) }
}
```

## Local http (debug only)

`app/src/debug/res/xml/network_security_config.xml`:

```xml
<network-security-config>
  <domain-config cleartextTrafficPermitted="true">
    <domain includeSubdomains="false">10.0.2.2</domain>
    <domain includeSubdomains="false">localhost</domain>
  </domain-config>
</network-security-config>
```

Reference it from `app/src/debug/AndroidManifest.xml` with `android:networkSecurityConfig="@xml/network_security_config"`. Make sure `<uses-permission android:name="android.permission.INTERNET"/>` exists in the main manifest.

## Paging

With Paging 3, use `offset` as the key: `nextKey = if (page.hasMore) page.offset + page.items.size else null`.

## Game menu (Jetpack Compose)

All three games must appear. Entries come from `/games`, with a built-in fallback.

```kotlin
private val fallbackGames = listOf("ftn" to "今彩539", "bl" to "大樂透", "pw" to "威力彩")

@Composable
fun GameMenu(games: List<GameInfo>?, selected: String, onSelect: (String) -> Unit) {
    val entries = games?.map { it.code to it.name } ?: fallbackGames
    NavigationBar {
        entries.forEach { (code, name) ->
            NavigationBarItem(
                selected = code == selected,
                onClick = { onSelect(code) },
                icon = { Icon(Icons.Outlined.Casino, contentDescription = null) },
                label = { Text(name) },
            )
        }
    }
}
```

Persist `selected` with DataStore, and route the three per-game screens by code (Navigation Compose routes `"{game}/draws"`, `"{game}/stats"`, `"{game}/analyze"`).
