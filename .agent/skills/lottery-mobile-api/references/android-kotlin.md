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
@Serializable
data class Envelope<T>(val code: Int, val message: String, val data: T? = null)

@Serializable
data class Features(
    val sum: Int, val odd: Int, val even: Int, val big: Int, val small: Int, val prime: Int,
    val tenGroup: List<Int>, val fullTenGroup: Boolean, val consecutive: String,
    val ac: Int, val span: Int, val mod3: List<Int>, val tailSum: Int,
    val sameTailGroups: Int, val maxGap: Int, val minGap: Int,
)

@Serializable
data class Draw(
    val period: String, val date: String, val yearIndex: Int,
    val numbers: List<Int>, val features: Features,
) { val localDate: LocalDate get() = LocalDate.parse(date) }

@Serializable
data class DrawPage(val items: List<Draw>, val total: Int, val limit: Int, val offset: Int) {
    val hasMore get() = offset + items.size < total
}

@Serializable
data class NumberStat(val number: Int, val windowCount: Int, val totalCount: Int, val omission: Int, val maxOmission: Int)

@Serializable
data class NumberStats(
    val window: Int, val totalDraws: Int, val expectedWindow: Double,
    val latestPeriod: String, val numbers: List<NumberStat>,
)

@Serializable data class AnalyzeRequest(val numbers: List<Int>)
@Serializable data class HistoryMatch(val match5: Int, val match4: Int, val match3: Int, val match2: Int)
@Serializable
data class AnalyzeResult(
    val numbers: List<Int>, val features: Features, val history: HistoryMatch,
    val totalDraws: Int, val lastMatch3: Draw? = null,
)

@Serializable
data class Health(
    val status: String, val draws: Int,
    val latestPeriod: String? = null, val latestDate: String? = null, val loadedAt: String? = null,
)
```

## Service and client

```kotlin
interface LotteryService {
    @GET("health") suspend fun health(): Response<Envelope<Health>>
    @GET("ftn/draws") suspend fun draws(
        @Query("limit") limit: Int = 20, @Query("offset") offset: Int = 0,
        @Query("from") from: String? = null, @Query("to") to: String? = null,
    ): Response<Envelope<DrawPage>>
    @GET("ftn/draws/latest") suspend fun latest(): Response<Envelope<Draw>>
    @GET("ftn/draws/{period}") suspend fun draw(@Path("period") period: String): Response<Envelope<Draw>>
    @GET("ftn/numbers/stats") suspend fun numberStats(@Query("window") window: Int = 100): Response<Envelope<NumberStats>>
    @POST("ftn/analyze") suspend fun analyze(@Body body: AnalyzeRequest): Response<Envelope<AnalyzeResult>>
}

class ApiException(val code: Int, override val message: String, val httpStatus: Int? = null) : Exception(message) {
    val retryable get() = code == -1 || code == 50301
}

object LotteryApi {
    private val json = Json { ignoreUnknownKeys = true }

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

fun validatePick(numbers: List<Int>) {
    require(numbers.size == 5 && numbers.toSet().size == 5 && numbers.all { it in 1..39 }) {
        "需選 5 個不重複的 1~39 號碼"
    }
}

// usage in a repository / ViewModel
suspend fun loadLatest(): Draw = LotteryApi.call { latest() }
suspend fun analyze(nums: List<Int>): AnalyzeResult {
    validatePick(nums)
    return LotteryApi.call { analyze(AnalyzeRequest(nums)) }
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

Reference it from `app/src/debug/AndroidManifest.xml` with `android:networkSecurityConfig="@xml/network_security_config"`. Also make sure `<uses-permission android:name="android.permission.INTERNET"/>` exists in the main manifest.

## Paging

With Paging 3, use `offset` as the key: `nextKey = if (page.hasMore) page.offset + page.items.size else null`.
