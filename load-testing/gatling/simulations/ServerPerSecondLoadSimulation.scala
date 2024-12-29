import io.gatling.core.Predef._
import io.gatling.http.Predef._
import scala.concurrent.duration._

class ServerPerSecondLoadSimulation extends Simulation {

  val projApi = http.baseUrl("http://172.28.22.1/api/v1")

  val scnEcho = scenario("Server Load Test")
    .exec(http("Get studio").get("/test/studios/1"))

  setUp(
    scnEcho.inject(
        constantUsersPerSec(2000).during(30)
    ).protocols(projApi),
  ).maxDuration(60.seconds)
}