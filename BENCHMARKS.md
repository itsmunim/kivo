# Benchmark Results

Tests run against kivo using `redis-benchmark`.

## Environment
- Date: 2026-09-13 15:15:33 UTC
- Host: Darwin MDS-MACBOOK-PRO.station 25.5.0 Darwin Kernel Version 25.5.0: Tue Jun  9 22:18:58 PDT 2026; root:xnu-12377.121.10~1/RELEASE_ARM64_T6000 arm64
- Go Version: go version go1.26.0 darwin/arm64

## Basic Operations (100000 requests, 50 parallel clients)

```
WARNING: Could not fetch server CONFIG
 PING_INLINE: rps=0.0 (overall: 0.0) avg_msec=0.000 (overall: 0.000)                                                                    PING_INLINE: rps=85400.0 (overall: 85400.0) avg_msec=0.342 (overall: 0.342)                                                                            PING_INLINE: rps=90151.4 (overall: 87780.4) avg_msec=0.304 (overall: 0.322)                                                                            PING_INLINE: rps=96552.0 (overall: 90700.4) avg_msec=0.287 (overall: 0.310)                                                                            PING_INLINE: rps=90900.0 (overall: 90750.2) avg_msec=0.298 (overall: 0.307)                                                                            ====== PING_INLINE ======
  100000 requests completed in 1.10 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.063 milliseconds (cumulative count 2)
50.000% <= 0.295 milliseconds (cumulative count 50870)
75.000% <= 0.319 milliseconds (cumulative count 77108)
87.500% <= 0.343 milliseconds (cumulative count 89037)
93.750% <= 0.367 milliseconds (cumulative count 93856)
96.875% <= 0.415 milliseconds (cumulative count 97013)
98.438% <= 0.559 milliseconds (cumulative count 98456)
99.219% <= 0.967 milliseconds (cumulative count 99220)
99.609% <= 1.151 milliseconds (cumulative count 99620)
99.805% <= 1.223 milliseconds (cumulative count 99813)
99.902% <= 1.663 milliseconds (cumulative count 99903)
99.951% <= 2.263 milliseconds (cumulative count 99954)
99.976% <= 2.359 milliseconds (cumulative count 99976)
99.988% <= 2.487 milliseconds (cumulative count 99988)
99.994% <= 2.503 milliseconds (cumulative count 99994)
99.997% <= 2.527 milliseconds (cumulative count 99998)
99.998% <= 2.543 milliseconds (cumulative count 99999)
99.999% <= 2.551 milliseconds (cumulative count 100000)
100.000% <= 2.551 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.016% <= 0.103 milliseconds (cumulative count 16)
3.049% <= 0.207 milliseconds (cumulative count 3049)
60.722% <= 0.303 milliseconds (cumulative count 60722)
96.729% <= 0.407 milliseconds (cumulative count 96729)
98.234% <= 0.503 milliseconds (cumulative count 98234)
98.562% <= 0.607 milliseconds (cumulative count 98562)
98.732% <= 0.703 milliseconds (cumulative count 98732)
99.044% <= 0.807 milliseconds (cumulative count 99044)
99.186% <= 0.903 milliseconds (cumulative count 99186)
99.271% <= 1.007 milliseconds (cumulative count 99271)
99.497% <= 1.103 milliseconds (cumulative count 99497)
99.782% <= 1.207 milliseconds (cumulative count 99782)
99.880% <= 1.303 milliseconds (cumulative count 99880)
99.883% <= 1.407 milliseconds (cumulative count 99883)
99.887% <= 1.503 milliseconds (cumulative count 99887)
99.894% <= 1.607 milliseconds (cumulative count 99894)
99.905% <= 1.703 milliseconds (cumulative count 99905)
99.911% <= 1.807 milliseconds (cumulative count 99911)
99.915% <= 1.903 milliseconds (cumulative count 99915)
99.922% <= 2.007 milliseconds (cumulative count 99922)
99.931% <= 2.103 milliseconds (cumulative count 99931)
100.000% <= 3.103 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 90991.81 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.306     0.056     0.295     0.383     0.799     2.551
 PING_MBULK: rps=54332.0 (overall: 91161.1) avg_msec=0.301 (overall: 0.301)                                                                           PING_MBULK: rps=89116.0 (overall: 89879.7) avg_msec=0.309 (overall: 0.306)                                                                           PING_MBULK: rps=90541.8 (overall: 90135.4) avg_msec=0.305 (overall: 0.306)                                                                           PING_MBULK: rps=91816.0 (overall: 90602.2) avg_msec=0.296 (overall: 0.303)                                                                           ====== PING_MBULK ======
  100000 requests completed in 1.10 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.079 milliseconds (cumulative count 5)
50.000% <= 0.303 milliseconds (cumulative count 57577)
75.000% <= 0.327 milliseconds (cumulative count 80847)
87.500% <= 0.343 milliseconds (cumulative count 88536)
93.750% <= 0.367 milliseconds (cumulative count 94257)
96.875% <= 0.399 milliseconds (cumulative count 97108)
98.438% <= 0.447 milliseconds (cumulative count 98587)
99.219% <= 0.511 milliseconds (cumulative count 99243)
99.609% <= 0.623 milliseconds (cumulative count 99624)
99.805% <= 0.743 milliseconds (cumulative count 99810)
99.902% <= 0.879 milliseconds (cumulative count 99903)
99.951% <= 1.151 milliseconds (cumulative count 99953)
99.976% <= 1.207 milliseconds (cumulative count 99976)
99.988% <= 1.263 milliseconds (cumulative count 99988)
99.994% <= 1.327 milliseconds (cumulative count 99994)
99.997% <= 1.383 milliseconds (cumulative count 99997)
99.998% <= 1.407 milliseconds (cumulative count 99999)
99.999% <= 1.415 milliseconds (cumulative count 100000)
100.000% <= 1.415 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.010% <= 0.103 milliseconds (cumulative count 10)
2.062% <= 0.207 milliseconds (cumulative count 2062)
57.577% <= 0.303 milliseconds (cumulative count 57577)
97.480% <= 0.407 milliseconds (cumulative count 97480)
99.204% <= 0.503 milliseconds (cumulative count 99204)
99.582% <= 0.607 milliseconds (cumulative count 99582)
99.746% <= 0.703 milliseconds (cumulative count 99746)
99.872% <= 0.807 milliseconds (cumulative count 99872)
99.911% <= 0.903 milliseconds (cumulative count 99911)
99.930% <= 1.007 milliseconds (cumulative count 99930)
99.940% <= 1.103 milliseconds (cumulative count 99940)
99.976% <= 1.207 milliseconds (cumulative count 99976)
99.992% <= 1.303 milliseconds (cumulative count 99992)
99.999% <= 1.407 milliseconds (cumulative count 99999)
100.000% <= 1.503 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 90991.81 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.301     0.072     0.303     0.375     0.479     1.415
 SET: rps=18924.0 (overall: 94620.0) avg_msec=0.308 (overall: 0.308)                                                                    SET: rps=92892.0 (overall: 93180.0) avg_msec=0.292 (overall: 0.295)                                                                    SET: rps=92312.0 (overall: 92785.5) avg_msec=0.294 (overall: 0.294)                                                                    SET: rps=90948.0 (overall: 92211.2) avg_msec=0.299 (overall: 0.296)                                                                    SET: rps=90000.0 (overall: 91683.2) avg_msec=0.302 (overall: 0.297)                                                                    ====== SET ======
  100000 requests completed in 1.09 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.063 milliseconds (cumulative count 1)
50.000% <= 0.303 milliseconds (cumulative count 59160)
75.000% <= 0.319 milliseconds (cumulative count 76824)
87.500% <= 0.343 milliseconds (cumulative count 90160)
93.750% <= 0.359 milliseconds (cumulative count 94159)
96.875% <= 0.383 milliseconds (cumulative count 97137)
98.438% <= 0.415 milliseconds (cumulative count 98654)
99.219% <= 0.455 milliseconds (cumulative count 99267)
99.609% <= 0.535 milliseconds (cumulative count 99621)
99.805% <= 0.615 milliseconds (cumulative count 99820)
99.902% <= 0.671 milliseconds (cumulative count 99911)
99.951% <= 0.727 milliseconds (cumulative count 99953)
99.976% <= 0.775 milliseconds (cumulative count 99977)
99.988% <= 0.855 milliseconds (cumulative count 99989)
99.994% <= 0.919 milliseconds (cumulative count 99994)
99.997% <= 0.943 milliseconds (cumulative count 99997)
99.998% <= 1.047 milliseconds (cumulative count 99999)
99.999% <= 1.175 milliseconds (cumulative count 100000)
100.000% <= 1.175 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.017% <= 0.103 milliseconds (cumulative count 17)
2.840% <= 0.207 milliseconds (cumulative count 2840)
59.160% <= 0.303 milliseconds (cumulative count 59160)
98.417% <= 0.407 milliseconds (cumulative count 98417)
99.525% <= 0.503 milliseconds (cumulative count 99525)
99.795% <= 0.607 milliseconds (cumulative count 99795)
99.936% <= 0.703 milliseconds (cumulative count 99936)
99.985% <= 0.807 milliseconds (cumulative count 99985)
99.992% <= 0.903 milliseconds (cumulative count 99992)
99.998% <= 1.007 milliseconds (cumulative count 99998)
99.999% <= 1.103 milliseconds (cumulative count 99999)
100.000% <= 1.207 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 91743.12 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.297     0.056     0.303     0.367     0.431     1.175
 GET: rps=77536.0 (overall: 92304.8) avg_msec=0.300 (overall: 0.300)                                                                    GET: rps=92736.0 (overall: 92539.1) avg_msec=0.294 (overall: 0.297)                                                                    GET: rps=91316.0 (overall: 92108.5) avg_msec=0.302 (overall: 0.299)                                                                    GET: rps=92768.9 (overall: 92281.0) avg_msec=0.292 (overall: 0.297)                                                                    ====== GET ======
  100000 requests completed in 1.09 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.047 milliseconds (cumulative count 1)
50.000% <= 0.295 milliseconds (cumulative count 50942)
75.000% <= 0.319 milliseconds (cumulative count 77488)
87.500% <= 0.343 milliseconds (cumulative count 89598)
93.750% <= 0.367 milliseconds (cumulative count 94749)
96.875% <= 0.391 milliseconds (cumulative count 97058)
98.438% <= 0.431 milliseconds (cumulative count 98542)
99.219% <= 0.495 milliseconds (cumulative count 99239)
99.609% <= 0.607 milliseconds (cumulative count 99630)
99.805% <= 0.695 milliseconds (cumulative count 99814)
99.902% <= 0.783 milliseconds (cumulative count 99904)
99.951% <= 0.967 milliseconds (cumulative count 99954)
99.976% <= 1.055 milliseconds (cumulative count 99977)
99.988% <= 1.111 milliseconds (cumulative count 99989)
99.994% <= 1.143 milliseconds (cumulative count 99994)
99.997% <= 2.175 milliseconds (cumulative count 99997)
99.998% <= 2.207 milliseconds (cumulative count 99999)
99.999% <= 2.223 milliseconds (cumulative count 100000)
100.000% <= 2.223 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.017% <= 0.103 milliseconds (cumulative count 17)
2.989% <= 0.207 milliseconds (cumulative count 2989)
61.177% <= 0.303 milliseconds (cumulative count 61177)
97.858% <= 0.407 milliseconds (cumulative count 97858)
99.292% <= 0.503 milliseconds (cumulative count 99292)
99.630% <= 0.607 milliseconds (cumulative count 99630)
99.827% <= 0.703 milliseconds (cumulative count 99827)
99.921% <= 0.807 milliseconds (cumulative count 99921)
99.942% <= 0.903 milliseconds (cumulative count 99942)
99.965% <= 1.007 milliseconds (cumulative count 99965)
99.987% <= 1.103 milliseconds (cumulative count 99987)
99.996% <= 1.207 milliseconds (cumulative count 99996)
100.000% <= 3.103 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 92165.90 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.297     0.040     0.295     0.375     0.471     2.223
 INCR: rps=45800.0 (overall: 91600.0) avg_msec=0.314 (overall: 0.314)                                                                     INCR: rps=94368.0 (overall: 93445.3) avg_msec=0.291 (overall: 0.298)                                                                     INCR: rps=93432.0 (overall: 93440.0) avg_msec=0.293 (overall: 0.296)                                                                     INCR: rps=92416.0 (overall: 93147.4) avg_msec=0.297 (overall: 0.297)                                                                     ====== INCR ======
  100000 requests completed in 1.08 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.063 milliseconds (cumulative count 1)
50.000% <= 0.295 milliseconds (cumulative count 52234)
75.000% <= 0.319 milliseconds (cumulative count 77907)
87.500% <= 0.343 milliseconds (cumulative count 89969)
93.750% <= 0.359 milliseconds (cumulative count 93795)
96.875% <= 0.391 milliseconds (cumulative count 97212)
98.438% <= 0.423 milliseconds (cumulative count 98597)
99.219% <= 0.463 milliseconds (cumulative count 99244)
99.609% <= 0.543 milliseconds (cumulative count 99625)
99.805% <= 0.687 milliseconds (cumulative count 99809)
99.902% <= 1.231 milliseconds (cumulative count 99903)
99.951% <= 1.695 milliseconds (cumulative count 99952)
99.976% <= 1.879 milliseconds (cumulative count 99977)
99.988% <= 1.927 milliseconds (cumulative count 99988)
99.994% <= 1.975 milliseconds (cumulative count 99994)
99.997% <= 2.095 milliseconds (cumulative count 99997)
99.998% <= 2.175 milliseconds (cumulative count 99999)
99.999% <= 2.215 milliseconds (cumulative count 100000)
100.000% <= 2.215 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.022% <= 0.103 milliseconds (cumulative count 22)
2.798% <= 0.207 milliseconds (cumulative count 2798)
61.980% <= 0.303 milliseconds (cumulative count 61980)
98.053% <= 0.407 milliseconds (cumulative count 98053)
99.493% <= 0.503 milliseconds (cumulative count 99493)
99.731% <= 0.607 milliseconds (cumulative count 99731)
99.823% <= 0.703 milliseconds (cumulative count 99823)
99.864% <= 0.807 milliseconds (cumulative count 99864)
99.878% <= 0.903 milliseconds (cumulative count 99878)
99.882% <= 1.007 milliseconds (cumulative count 99882)
99.891% <= 1.103 milliseconds (cumulative count 99891)
99.899% <= 1.207 milliseconds (cumulative count 99899)
99.907% <= 1.303 milliseconds (cumulative count 99907)
99.912% <= 1.407 milliseconds (cumulative count 99912)
99.919% <= 1.503 milliseconds (cumulative count 99919)
99.939% <= 1.607 milliseconds (cumulative count 99939)
99.952% <= 1.703 milliseconds (cumulative count 99952)
99.964% <= 1.807 milliseconds (cumulative count 99964)
99.983% <= 1.903 milliseconds (cumulative count 99983)
99.995% <= 2.007 milliseconds (cumulative count 99995)
99.997% <= 2.103 milliseconds (cumulative count 99997)
100.000% <= 3.103 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 92936.80 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.297     0.056     0.295     0.367     0.447     2.215
 LPUSH: rps=15123.5 (overall: 77469.4) avg_msec=0.532 (overall: 0.532)                                                                      LPUSH: rps=31340.0 (overall: 38899.7) avg_msec=1.563 (overall: 1.226)                                                                      LPUSH: rps=14184.0 (overall: 27644.8) avg_msec=3.501 (overall: 1.758)                                                                      LPUSH: rps=11016.0 (overall: 22441.8) avg_msec=4.527 (overall: 2.183)                                                                      LPUSH: rps=9740.0 (overall: 19414.7) avg_msec=5.115 (overall: 2.534)                                                                     LPUSH: rps=8520.0 (overall: 17317.9) avg_msec=5.875 (overall: 2.850)                                                                     LPUSH: rps=7556.0 (overall: 15742.4) avg_msec=6.580 (overall: 3.139)                                                                     LPUSH: rps=7168.0 (overall: 14550.9) avg_msec=6.981 (overall: 3.402)                                                                     LPUSH: rps=6645.4 (overall: 13582.9) avg_msec=7.468 (overall: 3.646)                                                                     LPUSH: rps=6192.0 (overall: 12779.6) avg_msec=8.056 (overall: 3.878)                                                                     LPUSH: rps=5728.0 (overall: 12088.2) avg_msec=8.726 (overall: 4.103)                                                                     LPUSH: rps=5628.0 (overall: 11511.4) avg_msec=8.869 (overall: 4.311)                                                                     LPUSH: rps=5876.0 (overall: 11049.5) avg_msec=8.530 (overall: 4.495)                                                                     LPUSH: rps=5712.0 (overall: 10645.2) avg_msec=8.753 (overall: 4.668)                                                                     LPUSH: rps=5484.0 (overall: 10281.7) avg_msec=9.094 (overall: 4.834)                                                                     LPUSH: rps=5432.0 (overall: 9962.6) avg_msec=9.211 (overall: 4.991)                                                                    LPUSH: rps=5364.0 (overall: 9678.8) avg_msec=9.275 (overall: 5.138)                                                                    LPUSH: rps=5055.8 (overall: 9409.0) avg_msec=9.811 (overall: 5.284)                                                                    LPUSH: rps=5148.0 (overall: 9174.9) avg_msec=9.762 (overall: 5.422)                                                                    LPUSH: rps=5356.0 (overall: 8976.0) avg_msec=9.341 (overall: 5.544)                                                                    LPUSH: rps=5212.0 (overall: 8789.7) avg_msec=9.545 (overall: 5.662)                                                                    LPUSH: rps=4756.0 (overall: 8599.5) avg_msec=10.513 (overall: 5.788)                                                                     LPUSH: rps=4608.0 (overall: 8419.7) avg_msec=10.848 (overall: 5.913)                                                                     LPUSH: rps=4430.3 (overall: 8247.2) avg_msec=11.249 (overall: 6.037)                                                                     LPUSH: rps=4340.0 (overall: 8085.8) avg_msec=11.480 (overall: 6.158)                                                                     LPUSH: rps=4052.0 (overall: 7925.7) avg_msec=12.330 (overall: 6.283)                                                                     LPUSH: rps=4116.0 (overall: 7780.4) avg_msec=12.129 (overall: 6.401)                                                                     LPUSH: rps=4060.0 (overall: 7643.6) avg_msec=12.331 (overall: 6.517)                                                                     LPUSH: rps=4147.4 (overall: 7519.2) avg_msec=12.041 (overall: 6.625)                                                                     LPUSH: rps=3968.0 (overall: 7397.6) avg_msec=12.604 (overall: 6.735)                                                                     LPUSH: rps=3652.0 (overall: 7273.7) avg_msec=13.627 (overall: 6.849)                                                                     LPUSH: rps=3716.0 (overall: 7159.7) avg_msec=13.507 (overall: 6.960)                                                                     LPUSH: rps=3760.0 (overall: 7054.1) avg_msec=13.260 (overall: 7.064)                                                                     LPUSH: rps=3864.0 (overall: 6958.1) avg_msec=12.908 (overall: 7.162)                                                                     LPUSH: rps=3684.0 (overall: 6862.4) avg_msec=13.605 (overall: 7.263)                                                                     LPUSH: rps=3462.2 (overall: 6765.4) avg_msec=14.356 (overall: 7.367)                                                                     LPUSH: rps=3512.0 (overall: 6675.6) avg_msec=14.217 (overall: 7.466)                                                                     LPUSH: rps=3444.0 (overall: 6588.8) avg_msec=14.518 (overall: 7.565)                                                                     LPUSH: rps=3388.0 (overall: 6505.0) avg_msec=14.747 (overall: 7.663)                                                                     LPUSH: rps=3312.0 (overall: 6423.6) avg_msec=15.149 (overall: 7.761)                                                                     LPUSH: rps=3280.0 (overall: 6345.4) avg_msec=15.180 (overall: 7.857)                                                                     LPUSH: rps=3286.9 (overall: 6270.9) avg_msec=15.165 (overall: 7.950)                                                                     LPUSH: rps=3300.0 (overall: 6200.6) avg_msec=15.080 (overall: 8.040)                                                                     LPUSH: rps=3304.0 (overall: 6133.5) avg_msec=15.191 (overall: 8.129)                                                                     LPUSH: rps=3284.0 (overall: 6069.1) avg_msec=15.224 (overall: 8.216)                                                                     LPUSH: rps=3228.0 (overall: 6006.3) avg_msec=15.457 (overall: 8.302)                                                                     LPUSH: rps=3168.0 (overall: 5944.9) avg_msec=15.794 (overall: 8.388)                                                                     LPUSH: rps=3156.0 (overall: 5885.8) avg_msec=15.837 (overall: 8.473)                                                                     LPUSH: rps=3015.9 (overall: 5826.1) avg_msec=16.465 (overall: 8.559)                                                                     LPUSH: rps=3004.0 (overall: 5768.7) avg_msec=16.610 (overall: 8.644)                                                                     LPUSH: rps=2804.0 (overall: 5709.7) avg_msec=17.849 (overall: 8.734)                                                                     LPUSH: rps=2772.0 (overall: 5652.4) avg_msec=17.978 (overall: 8.823)                                                                     LPUSH: rps=2876.0 (overall: 5599.2) avg_msec=17.486 (overall: 8.908)                                                                     LPUSH: rps=2720.0 (overall: 5545.1) avg_msec=18.179 (overall: 8.993)                                                                     LPUSH: rps=2840.6 (overall: 5495.0) avg_msec=17.745 (overall: 9.077)                                                                     LPUSH: rps=2912.0 (overall: 5448.3) avg_msec=17.131 (overall: 9.155)                                                                     LPUSH: rps=2760.0 (overall: 5400.4) avg_msec=18.112 (overall: 9.237)                                                                     LPUSH: rps=2800.0 (overall: 5355.0) avg_msec=17.865 (overall: 9.315)                                                                     LPUSH: rps=2748.0 (overall: 5310.2) avg_msec=18.112 (overall: 9.394)                                                                     LPUSH: rps=2656.0 (overall: 5265.4) avg_msec=18.692 (overall: 9.473)                                                                     LPUSH: rps=2717.1 (overall: 5222.9) avg_msec=18.429 (overall: 9.550)                                                                     LPUSH: rps=2740.0 (overall: 5182.4) avg_msec=18.396 (overall: 9.627)                                                                     LPUSH: rps=2696.0 (overall: 5142.4) avg_msec=18.542 (overall: 9.702)                                                                     LPUSH: rps=2568.0 (overall: 5101.7) avg_msec=19.388 (overall: 9.779)                                                                     LPUSH: rps=2464.0 (overall: 5060.7) avg_msec=20.292 (overall: 9.859)                                                                     LPUSH: rps=2528.0 (overall: 5021.8) avg_msec=19.693 (overall: 9.935)                                                                     LPUSH: rps=2648.0 (overall: 4986.0) avg_msec=18.931 (overall: 10.007)                                                                      LPUSH: rps=2585.7 (overall: 4950.1) avg_msec=19.303 (overall: 10.079)                                                                      LPUSH: rps=2640.0 (overall: 4916.3) avg_msec=18.934 (overall: 10.149)                                                                      LPUSH: rps=2704.0 (overall: 4884.3) avg_msec=18.505 (overall: 10.216)                                                                      LPUSH: rps=2648.0 (overall: 4852.5) avg_msec=18.902 (overall: 10.283)                                                                      LPUSH: rps=2540.0 (overall: 4820.0) avg_msec=19.636 (overall: 10.352)                                                                      LPUSH: rps=2468.0 (overall: 4787.5) avg_msec=20.093 (overall: 10.422)                                                                      LPUSH: rps=2506.0 (overall: 4756.2) avg_msec=19.913 (overall: 10.490)                                                                      LPUSH: rps=2528.0 (overall: 4726.2) avg_msec=19.938 (overall: 10.559)                                                                      LPUSH: rps=2536.0 (overall: 4697.1) avg_msec=19.695 (overall: 10.624)                                                                      LPUSH: rps=2444.0 (overall: 4667.5) avg_msec=20.426 (overall: 10.691)                                                                      LPUSH: rps=2512.0 (overall: 4639.6) avg_msec=19.867 (overall: 10.756)                                                                      LPUSH: rps=2456.0 (overall: 4611.7) avg_msec=20.302 (overall: 10.821)                                                                      LPUSH: rps=2366.5 (overall: 4583.3) avg_msec=21.153 (overall: 10.888)                                                                      LPUSH: rps=2376.0 (overall: 4555.8) avg_msec=20.987 (overall: 10.954)                                                                      LPUSH: rps=2296.0 (overall: 4527.9) avg_msec=21.728 (overall: 11.021)                                                                      LPUSH: rps=2376.0 (overall: 4501.8) avg_msec=21.037 (overall: 11.085)                                                                      LPUSH: rps=2388.0 (overall: 4476.4) avg_msec=21.047 (overall: 11.149)                                                                      LPUSH: rps=2215.1 (overall: 4449.4) avg_msec=22.333 (overall: 11.216)                                                                      LPUSH: rps=2380.0 (overall: 4425.2) avg_msec=21.024 (overall: 11.278)                                                                      LPUSH: rps=2272.0 (overall: 4400.2) avg_msec=22.020 (overall: 11.342)                                                                      LPUSH: rps=2336.0 (overall: 4376.5) avg_msec=21.624 (overall: 11.405)                                                                      LPUSH: rps=2308.0 (overall: 4353.1) avg_msec=21.497 (overall: 11.465)                                                                      LPUSH: rps=2290.8 (overall: 4329.9) avg_msec=21.752 (overall: 11.527)                                                                      LPUSH: rps=2220.0 (overall: 4306.5) avg_msec=22.373 (overall: 11.589)                                                                      LPUSH: rps=2248.0 (overall: 4284.0) avg_msec=22.432 (overall: 11.651)                                                                      LPUSH: rps=2240.0 (overall: 4261.8) avg_msec=22.259 (overall: 11.711)                                                                      LPUSH: rps=2270.9 (overall: 4240.4) avg_msec=21.967 (overall: 11.770)                                                                      LPUSH: rps=2192.0 (overall: 4218.6) avg_msec=22.818 (overall: 11.831)                                                                      LPUSH: rps=2208.0 (overall: 4197.5) avg_msec=22.619 (overall: 11.891)                                                                      ====== LPUSH ======
  100000 requests completed in 23.83 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.023 milliseconds (cumulative count 4)
50.000% <= 11.839 milliseconds (cumulative count 50008)
75.000% <= 17.663 milliseconds (cumulative count 75053)
87.500% <= 20.223 milliseconds (cumulative count 87545)
93.750% <= 21.711 milliseconds (cumulative count 93750)
96.875% <= 22.671 milliseconds (cumulative count 96894)
98.438% <= 23.359 milliseconds (cumulative count 98451)
99.219% <= 23.983 milliseconds (cumulative count 99224)
99.609% <= 24.511 milliseconds (cumulative count 99611)
99.805% <= 25.039 milliseconds (cumulative count 99805)
99.902% <= 25.775 milliseconds (cumulative count 99904)
99.951% <= 28.367 milliseconds (cumulative count 99952)
99.976% <= 37.055 milliseconds (cumulative count 99976)
99.988% <= 41.535 milliseconds (cumulative count 99989)
99.994% <= 43.647 milliseconds (cumulative count 99994)
99.997% <= 44.895 milliseconds (cumulative count 99997)
99.998% <= 45.599 milliseconds (cumulative count 99999)
99.999% <= 46.175 milliseconds (cumulative count 100000)
100.000% <= 46.175 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.987% <= 0.103 milliseconds (cumulative count 987)
1.764% <= 0.207 milliseconds (cumulative count 1764)
2.417% <= 0.303 milliseconds (cumulative count 2417)
3.296% <= 0.407 milliseconds (cumulative count 3296)
3.532% <= 0.503 milliseconds (cumulative count 3532)
3.750% <= 0.607 milliseconds (cumulative count 3750)
3.976% <= 0.703 milliseconds (cumulative count 3976)
4.269% <= 0.807 milliseconds (cumulative count 4269)
4.606% <= 0.903 milliseconds (cumulative count 4606)
4.970% <= 1.007 milliseconds (cumulative count 4970)
5.475% <= 1.103 milliseconds (cumulative count 5475)
6.351% <= 1.207 milliseconds (cumulative count 6351)
7.121% <= 1.303 milliseconds (cumulative count 7121)
7.572% <= 1.407 milliseconds (cumulative count 7572)
7.930% <= 1.503 milliseconds (cumulative count 7930)
8.349% <= 1.607 milliseconds (cumulative count 8349)
8.707% <= 1.703 milliseconds (cumulative count 8707)
9.107% <= 1.807 milliseconds (cumulative count 9107)
9.468% <= 1.903 milliseconds (cumulative count 9468)
9.680% <= 2.007 milliseconds (cumulative count 9680)
9.873% <= 2.103 milliseconds (cumulative count 9873)
12.460% <= 3.103 milliseconds (cumulative count 12460)
15.363% <= 4.103 milliseconds (cumulative count 15363)
19.158% <= 5.103 milliseconds (cumulative count 19158)
22.429% <= 6.103 milliseconds (cumulative count 22429)
25.869% <= 7.103 milliseconds (cumulative count 25869)
29.878% <= 8.103 milliseconds (cumulative count 29878)
36.973% <= 9.103 milliseconds (cumulative count 36973)
43.396% <= 10.103 milliseconds (cumulative count 43396)
47.173% <= 11.103 milliseconds (cumulative count 47173)
50.994% <= 12.103 milliseconds (cumulative count 50994)
55.169% <= 13.103 milliseconds (cumulative count 55169)
59.687% <= 14.103 milliseconds (cumulative count 59687)
64.423% <= 15.103 milliseconds (cumulative count 64423)
68.598% <= 16.103 milliseconds (cumulative count 68598)
72.620% <= 17.103 milliseconds (cumulative count 72620)
77.140% <= 18.111 milliseconds (cumulative count 77140)
82.278% <= 19.103 milliseconds (cumulative count 82278)
87.032% <= 20.111 milliseconds (cumulative count 87032)
91.424% <= 21.103 milliseconds (cumulative count 91424)
95.166% <= 22.111 milliseconds (cumulative count 95166)
97.963% <= 23.103 milliseconds (cumulative count 97963)
99.333% <= 24.111 milliseconds (cumulative count 99333)
99.825% <= 25.103 milliseconds (cumulative count 99825)
99.917% <= 26.111 milliseconds (cumulative count 99917)
99.938% <= 27.103 milliseconds (cumulative count 99938)
99.949% <= 28.111 milliseconds (cumulative count 99949)
99.958% <= 29.103 milliseconds (cumulative count 99958)
99.961% <= 30.111 milliseconds (cumulative count 99961)
99.962% <= 31.103 milliseconds (cumulative count 99962)
99.967% <= 32.111 milliseconds (cumulative count 99967)
99.968% <= 33.119 milliseconds (cumulative count 99968)
99.970% <= 34.111 milliseconds (cumulative count 99970)
99.971% <= 35.103 milliseconds (cumulative count 99971)
99.974% <= 36.127 milliseconds (cumulative count 99974)
99.976% <= 37.119 milliseconds (cumulative count 99976)
99.981% <= 39.103 milliseconds (cumulative count 99981)
99.983% <= 40.127 milliseconds (cumulative count 99983)
99.986% <= 41.119 milliseconds (cumulative count 99986)
99.990% <= 42.111 milliseconds (cumulative count 99990)
99.992% <= 43.103 milliseconds (cumulative count 99992)
99.996% <= 44.127 milliseconds (cumulative count 99996)
99.998% <= 45.119 milliseconds (cumulative count 99998)
99.999% <= 46.111 milliseconds (cumulative count 99999)
100.000% <= 47.103 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 4196.04 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
       11.895     0.016    11.839    22.063    23.775    46.175
 RPUSH: rps=80436.0 (overall: 87051.9) avg_msec=0.342 (overall: 0.342)                                                                      RPUSH: rps=93948.0 (overall: 90636.2) avg_msec=0.293 (overall: 0.315)                                                                      RPUSH: rps=92629.5 (overall: 91319.7) avg_msec=0.296 (overall: 0.309)                                                                      RPUSH: rps=96280.0 (overall: 92582.5) avg_msec=0.286 (overall: 0.303)                                                                      ====== RPUSH ======
  100000 requests completed in 1.08 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.047 milliseconds (cumulative count 1)
50.000% <= 0.295 milliseconds (cumulative count 55962)
75.000% <= 0.319 milliseconds (cumulative count 80323)
87.500% <= 0.335 milliseconds (cumulative count 88460)
93.750% <= 0.359 milliseconds (cumulative count 94294)
96.875% <= 0.391 milliseconds (cumulative count 96955)
98.438% <= 0.647 milliseconds (cumulative count 98440)
99.219% <= 1.095 milliseconds (cumulative count 99225)
99.609% <= 1.215 milliseconds (cumulative count 99611)
99.805% <= 1.303 milliseconds (cumulative count 99806)
99.902% <= 2.039 milliseconds (cumulative count 99903)
99.951% <= 2.175 milliseconds (cumulative count 99952)
99.976% <= 2.247 milliseconds (cumulative count 99976)
99.988% <= 2.375 milliseconds (cumulative count 99988)
99.994% <= 2.471 milliseconds (cumulative count 99995)
99.997% <= 2.487 milliseconds (cumulative count 99998)
99.998% <= 2.503 milliseconds (cumulative count 99999)
99.999% <= 2.575 milliseconds (cumulative count 100000)
100.000% <= 2.575 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.018% <= 0.103 milliseconds (cumulative count 18)
3.153% <= 0.207 milliseconds (cumulative count 3153)
65.749% <= 0.303 milliseconds (cumulative count 65749)
97.495% <= 0.407 milliseconds (cumulative count 97495)
98.279% <= 0.503 milliseconds (cumulative count 98279)
98.423% <= 0.607 milliseconds (cumulative count 98423)
98.484% <= 0.703 milliseconds (cumulative count 98484)
98.776% <= 0.807 milliseconds (cumulative count 98776)
99.033% <= 0.903 milliseconds (cumulative count 99033)
99.063% <= 1.007 milliseconds (cumulative count 99063)
99.234% <= 1.103 milliseconds (cumulative count 99234)
99.594% <= 1.207 milliseconds (cumulative count 99594)
99.806% <= 1.303 milliseconds (cumulative count 99806)
99.849% <= 1.407 milliseconds (cumulative count 99849)
99.850% <= 1.503 milliseconds (cumulative count 99850)
99.880% <= 1.607 milliseconds (cumulative count 99880)
99.895% <= 1.703 milliseconds (cumulative count 99895)
99.901% <= 1.807 milliseconds (cumulative count 99901)
99.902% <= 1.903 milliseconds (cumulative count 99902)
99.915% <= 2.103 milliseconds (cumulative count 99915)
100.000% <= 3.103 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 92421.44 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.303     0.040     0.295     0.367     0.863     2.575
 LPOP: rps=55532.0 (overall: 93174.5) avg_msec=0.298 (overall: 0.298)                                                                     LPOP: rps=89496.0 (overall: 90869.7) avg_msec=0.309 (overall: 0.305)                                                                     LPOP: rps=90657.4 (overall: 90787.7) avg_msec=0.301 (overall: 0.303)                                                                     LPOP: rps=89144.0 (overall: 90331.1) avg_msec=0.303 (overall: 0.303)                                                                     ====== LPOP ======
  100000 requests completed in 1.11 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.071 milliseconds (cumulative count 2)
50.000% <= 0.303 milliseconds (cumulative count 56049)
75.000% <= 0.319 milliseconds (cumulative count 76066)
87.500% <= 0.343 milliseconds (cumulative count 89416)
93.750% <= 0.367 milliseconds (cumulative count 94345)
96.875% <= 0.399 milliseconds (cumulative count 97094)
98.438% <= 0.439 milliseconds (cumulative count 98567)
99.219% <= 0.511 milliseconds (cumulative count 99227)
99.609% <= 0.639 milliseconds (cumulative count 99624)
99.805% <= 0.751 milliseconds (cumulative count 99809)
99.902% <= 0.847 milliseconds (cumulative count 99903)
99.951% <= 0.959 milliseconds (cumulative count 99954)
99.976% <= 1.087 milliseconds (cumulative count 99976)
99.988% <= 1.271 milliseconds (cumulative count 99989)
99.994% <= 1.399 milliseconds (cumulative count 99994)
99.997% <= 1.439 milliseconds (cumulative count 99997)
99.998% <= 1.455 milliseconds (cumulative count 99999)
99.999% <= 1.671 milliseconds (cumulative count 100000)
100.000% <= 1.671 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.011% <= 0.103 milliseconds (cumulative count 11)
1.608% <= 0.207 milliseconds (cumulative count 1608)
56.049% <= 0.303 milliseconds (cumulative count 56049)
97.552% <= 0.407 milliseconds (cumulative count 97552)
99.190% <= 0.503 milliseconds (cumulative count 99190)
99.553% <= 0.607 milliseconds (cumulative count 99553)
99.750% <= 0.703 milliseconds (cumulative count 99750)
99.874% <= 0.807 milliseconds (cumulative count 99874)
99.933% <= 0.903 milliseconds (cumulative count 99933)
99.962% <= 1.007 milliseconds (cumulative count 99962)
99.976% <= 1.103 milliseconds (cumulative count 99976)
99.985% <= 1.207 milliseconds (cumulative count 99985)
99.990% <= 1.303 milliseconds (cumulative count 99990)
99.995% <= 1.407 milliseconds (cumulative count 99995)
99.999% <= 1.503 milliseconds (cumulative count 99999)
100.000% <= 1.703 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 90171.33 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.303     0.064     0.303     0.375     0.479     1.671
 RPOP: rps=15108.0 (overall: 94425.0) avg_msec=0.314 (overall: 0.314)                                                                     RPOP: rps=93486.1 (overall: 93615.1) avg_msec=0.294 (overall: 0.297)                                                                     RPOP: rps=94336.0 (overall: 93948.2) avg_msec=0.291 (overall: 0.294)                                                                     RPOP: rps=93260.0 (overall: 93730.7) avg_msec=0.294 (overall: 0.294)                                                                     RPOP: rps=93916.0 (overall: 93775.2) avg_msec=0.291 (overall: 0.293)                                                                     ====== RPOP ======
  100000 requests completed in 1.07 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.055 milliseconds (cumulative count 1)
50.000% <= 0.295 milliseconds (cumulative count 55248)
75.000% <= 0.319 milliseconds (cumulative count 79489)
87.500% <= 0.335 milliseconds (cumulative count 88019)
93.750% <= 0.359 milliseconds (cumulative count 94593)
96.875% <= 0.383 milliseconds (cumulative count 97471)
98.438% <= 0.407 milliseconds (cumulative count 98645)
99.219% <= 0.431 milliseconds (cumulative count 99228)
99.609% <= 0.471 milliseconds (cumulative count 99618)
99.805% <= 0.535 milliseconds (cumulative count 99810)
99.902% <= 0.671 milliseconds (cumulative count 99904)
99.951% <= 0.783 milliseconds (cumulative count 99956)
99.976% <= 0.871 milliseconds (cumulative count 99976)
99.988% <= 0.959 milliseconds (cumulative count 99988)
99.994% <= 0.991 milliseconds (cumulative count 99994)
99.997% <= 1.039 milliseconds (cumulative count 99998)
99.998% <= 1.047 milliseconds (cumulative count 99999)
99.999% <= 1.055 milliseconds (cumulative count 100000)
100.000% <= 1.055 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.012% <= 0.103 milliseconds (cumulative count 12)
2.539% <= 0.207 milliseconds (cumulative count 2539)
64.590% <= 0.303 milliseconds (cumulative count 64590)
98.645% <= 0.407 milliseconds (cumulative count 98645)
99.756% <= 0.503 milliseconds (cumulative count 99756)
99.867% <= 0.607 milliseconds (cumulative count 99867)
99.925% <= 0.703 milliseconds (cumulative count 99925)
99.962% <= 0.807 milliseconds (cumulative count 99962)
99.982% <= 0.903 milliseconds (cumulative count 99982)
99.994% <= 1.007 milliseconds (cumulative count 99994)
100.000% <= 1.103 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 93808.63 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.293     0.048     0.295     0.367     0.423     1.055
 SADD: rps=84115.5 (overall: 93835.6) avg_msec=0.295 (overall: 0.295)                                                                     SADD: rps=92708.0 (overall: 93242.1) avg_msec=0.293 (overall: 0.294)                                                                     SADD: rps=92332.0 (overall: 92928.3) avg_msec=0.298 (overall: 0.296)                                                                     SADD: rps=97688.0 (overall: 94148.7) avg_msec=0.285 (overall: 0.293)                                                                     ====== SADD ======
  100000 requests completed in 1.07 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.055 milliseconds (cumulative count 1)
50.000% <= 0.295 milliseconds (cumulative count 53423)
75.000% <= 0.319 milliseconds (cumulative count 79158)
87.500% <= 0.335 milliseconds (cumulative count 87765)
93.750% <= 0.367 milliseconds (cumulative count 94881)
96.875% <= 0.391 milliseconds (cumulative count 96952)
98.438% <= 0.423 milliseconds (cumulative count 98486)
99.219% <= 0.479 milliseconds (cumulative count 99258)
99.609% <= 0.599 milliseconds (cumulative count 99619)
99.805% <= 0.727 milliseconds (cumulative count 99811)
99.902% <= 0.831 milliseconds (cumulative count 99906)
99.951% <= 0.991 milliseconds (cumulative count 99952)
99.976% <= 1.439 milliseconds (cumulative count 99976)
99.988% <= 1.607 milliseconds (cumulative count 99988)
99.994% <= 1.679 milliseconds (cumulative count 99994)
99.997% <= 1.767 milliseconds (cumulative count 99997)
99.998% <= 1.823 milliseconds (cumulative count 99999)
99.999% <= 1.847 milliseconds (cumulative count 100000)
100.000% <= 1.847 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.023% <= 0.103 milliseconds (cumulative count 23)
4.154% <= 0.207 milliseconds (cumulative count 4154)
63.302% <= 0.303 milliseconds (cumulative count 63302)
97.862% <= 0.407 milliseconds (cumulative count 97862)
99.374% <= 0.503 milliseconds (cumulative count 99374)
99.637% <= 0.607 milliseconds (cumulative count 99637)
99.774% <= 0.703 milliseconds (cumulative count 99774)
99.891% <= 0.807 milliseconds (cumulative count 99891)
99.940% <= 0.903 milliseconds (cumulative count 99940)
99.952% <= 1.007 milliseconds (cumulative count 99952)
99.956% <= 1.103 milliseconds (cumulative count 99956)
99.959% <= 1.207 milliseconds (cumulative count 99959)
99.971% <= 1.303 milliseconds (cumulative count 99971)
99.974% <= 1.407 milliseconds (cumulative count 99974)
99.983% <= 1.503 milliseconds (cumulative count 99983)
99.988% <= 1.607 milliseconds (cumulative count 99988)
99.994% <= 1.703 milliseconds (cumulative count 99994)
99.998% <= 1.807 milliseconds (cumulative count 99998)
100.000% <= 1.903 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 93720.71 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.294     0.048     0.295     0.375     0.455     1.847
 HSET: rps=58772.0 (overall: 94185.9) avg_msec=0.300 (overall: 0.300)                                                                     HSET: rps=95211.2 (overall: 94818.2) avg_msec=0.288 (overall: 0.292)                                                                     HSET: rps=91556.0 (overall: 93576.9) avg_msec=0.301 (overall: 0.296)                                                                     HSET: rps=87812.0 (overall: 91987.9) avg_msec=0.323 (overall: 0.303)                                                                     ====== HSET ======
  100000 requests completed in 1.10 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.039 milliseconds (cumulative count 2)
50.000% <= 0.303 milliseconds (cumulative count 57103)
75.000% <= 0.327 milliseconds (cumulative count 78178)
87.500% <= 0.351 milliseconds (cumulative count 88213)
93.750% <= 0.391 milliseconds (cumulative count 94520)
96.875% <= 0.439 milliseconds (cumulative count 96932)
98.438% <= 0.551 milliseconds (cumulative count 98495)
99.219% <= 0.703 milliseconds (cumulative count 99243)
99.609% <= 0.919 milliseconds (cumulative count 99616)
99.805% <= 1.135 milliseconds (cumulative count 99809)
99.902% <= 1.567 milliseconds (cumulative count 99906)
99.951% <= 2.503 milliseconds (cumulative count 99952)
99.976% <= 2.607 milliseconds (cumulative count 99978)
99.988% <= 2.655 milliseconds (cumulative count 99988)
99.994% <= 2.759 milliseconds (cumulative count 99994)
99.997% <= 2.831 milliseconds (cumulative count 99997)
99.998% <= 2.863 milliseconds (cumulative count 99999)
99.999% <= 2.879 milliseconds (cumulative count 100000)
100.000% <= 2.879 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.067% <= 0.103 milliseconds (cumulative count 67)
2.682% <= 0.207 milliseconds (cumulative count 2682)
57.103% <= 0.303 milliseconds (cumulative count 57103)
95.653% <= 0.407 milliseconds (cumulative count 95653)
98.039% <= 0.503 milliseconds (cumulative count 98039)
98.842% <= 0.607 milliseconds (cumulative count 98842)
99.243% <= 0.703 milliseconds (cumulative count 99243)
99.492% <= 0.807 milliseconds (cumulative count 99492)
99.601% <= 0.903 milliseconds (cumulative count 99601)
99.697% <= 1.007 milliseconds (cumulative count 99697)
99.782% <= 1.103 milliseconds (cumulative count 99782)
99.838% <= 1.207 milliseconds (cumulative count 99838)
99.862% <= 1.303 milliseconds (cumulative count 99862)
99.874% <= 1.407 milliseconds (cumulative count 99874)
99.891% <= 1.503 milliseconds (cumulative count 99891)
99.917% <= 1.607 milliseconds (cumulative count 99917)
99.939% <= 1.703 milliseconds (cumulative count 99939)
99.950% <= 1.807 milliseconds (cumulative count 99950)
100.000% <= 3.103 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 90579.71 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.308     0.032     0.303     0.399     0.639     2.879

```

## String Operations (100000 requests, 50 parallel clients)

```
WARNING: Could not fetch server CONFIG
 SET: rps=0.0 (overall: 0.0) avg_msec=0.000 (overall: 0.000)                                                            SET: rps=86340.0 (overall: 85996.0) avg_msec=0.325 (overall: 0.325)                                                                    SET: rps=82956.0 (overall: 84479.0) avg_msec=0.340 (overall: 0.332)                                                                    SET: rps=91880.0 (overall: 86942.7) avg_msec=0.300 (overall: 0.321)                                                                    SET: rps=90432.0 (overall: 87814.2) avg_msec=0.301 (overall: 0.316)                                                                    ====== SET ======
  100000 requests completed in 1.13 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.031 milliseconds (cumulative count 1)
50.000% <= 0.311 milliseconds (cumulative count 59102)
75.000% <= 0.335 milliseconds (cumulative count 79608)
87.500% <= 0.359 milliseconds (cumulative count 89043)
93.750% <= 0.391 milliseconds (cumulative count 94396)
96.875% <= 0.431 milliseconds (cumulative count 96966)
98.438% <= 0.511 milliseconds (cumulative count 98451)
99.219% <= 0.679 milliseconds (cumulative count 99225)
99.609% <= 0.895 milliseconds (cumulative count 99612)
99.805% <= 1.079 milliseconds (cumulative count 99810)
99.902% <= 1.271 milliseconds (cumulative count 99905)
99.951% <= 1.407 milliseconds (cumulative count 99955)
99.976% <= 1.495 milliseconds (cumulative count 99977)
99.988% <= 1.639 milliseconds (cumulative count 99988)
99.994% <= 1.711 milliseconds (cumulative count 99994)
99.997% <= 1.743 milliseconds (cumulative count 99997)
99.998% <= 1.767 milliseconds (cumulative count 99999)
99.999% <= 2.527 milliseconds (cumulative count 100000)
100.000% <= 2.527 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.029% <= 0.103 milliseconds (cumulative count 29)
1.688% <= 0.207 milliseconds (cumulative count 1688)
49.413% <= 0.303 milliseconds (cumulative count 49413)
95.700% <= 0.407 milliseconds (cumulative count 95700)
98.375% <= 0.503 milliseconds (cumulative count 98375)
99.016% <= 0.607 milliseconds (cumulative count 99016)
99.283% <= 0.703 milliseconds (cumulative count 99283)
99.489% <= 0.807 milliseconds (cumulative count 99489)
99.622% <= 0.903 milliseconds (cumulative count 99622)
99.744% <= 1.007 milliseconds (cumulative count 99744)
99.826% <= 1.103 milliseconds (cumulative count 99826)
99.877% <= 1.207 milliseconds (cumulative count 99877)
99.917% <= 1.303 milliseconds (cumulative count 99917)
99.955% <= 1.407 milliseconds (cumulative count 99955)
99.979% <= 1.503 milliseconds (cumulative count 99979)
99.984% <= 1.607 milliseconds (cumulative count 99984)
99.992% <= 1.703 milliseconds (cumulative count 99992)
99.999% <= 1.807 milliseconds (cumulative count 99999)
100.000% <= 3.103 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 88417.33 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.313     0.024     0.311     0.399     0.607     2.527
 GET: rps=43135.5 (overall: 90225.0) avg_msec=0.315 (overall: 0.315)                                                                    GET: rps=91776.0 (overall: 91273.0) avg_msec=0.301 (overall: 0.305)                                                                    GET: rps=93228.0 (overall: 92061.3) avg_msec=0.295 (overall: 0.301)                                                                    GET: rps=92636.0 (overall: 92226.4) avg_msec=0.296 (overall: 0.300)                                                                    ====== GET ======
  100000 requests completed in 1.08 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.039 milliseconds (cumulative count 2)
50.000% <= 0.295 milliseconds (cumulative count 51100)
75.000% <= 0.319 milliseconds (cumulative count 76677)
87.500% <= 0.343 milliseconds (cumulative count 88788)
93.750% <= 0.367 milliseconds (cumulative count 94205)
96.875% <= 0.399 milliseconds (cumulative count 97322)
98.438% <= 0.431 milliseconds (cumulative count 98491)
99.219% <= 0.495 milliseconds (cumulative count 99222)
99.609% <= 0.615 milliseconds (cumulative count 99615)
99.805% <= 0.743 milliseconds (cumulative count 99813)
99.902% <= 0.823 milliseconds (cumulative count 99914)
99.951% <= 0.895 milliseconds (cumulative count 99959)
99.976% <= 0.959 milliseconds (cumulative count 99978)
99.988% <= 0.999 milliseconds (cumulative count 99991)
99.994% <= 1.031 milliseconds (cumulative count 99994)
99.997% <= 1.143 milliseconds (cumulative count 99997)
99.998% <= 1.175 milliseconds (cumulative count 99999)
99.999% <= 1.351 milliseconds (cumulative count 100000)
100.000% <= 1.351 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.041% <= 0.103 milliseconds (cumulative count 41)
2.227% <= 0.207 milliseconds (cumulative count 2227)
61.135% <= 0.303 milliseconds (cumulative count 61135)
97.711% <= 0.407 milliseconds (cumulative count 97711)
99.269% <= 0.503 milliseconds (cumulative count 99269)
99.601% <= 0.607 milliseconds (cumulative count 99601)
99.761% <= 0.703 milliseconds (cumulative count 99761)
99.886% <= 0.807 milliseconds (cumulative count 99886)
99.961% <= 0.903 milliseconds (cumulative count 99961)
99.991% <= 1.007 milliseconds (cumulative count 99991)
99.996% <= 1.103 milliseconds (cumulative count 99996)
99.999% <= 1.207 milliseconds (cumulative count 99999)
100.000% <= 1.407 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 92421.44 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.299     0.032     0.295     0.375     0.471     1.351

```

## List Operations (100000 requests, 50 parallel clients)

```
WARNING: Could not fetch server CONFIG
 LPUSH: rps=0.0 (overall: 0.0) avg_msec=0.000 (overall: 0.000)                                                              LPUSH: rps=43340.0 (overall: 43340.0) avg_msec=1.098 (overall: 1.098)                                                                      LPUSH: rps=15604.0 (overall: 29472.0) avg_msec=3.184 (overall: 1.650)                                                                      LPUSH: rps=11364.0 (overall: 23436.0) avg_msec=4.386 (overall: 2.092)                                                                      LPUSH: rps=9968.0 (overall: 20069.0) avg_msec=5.001 (overall: 2.454)                                                                     LPUSH: rps=8576.0 (overall: 17770.4) avg_msec=5.818 (overall: 2.778)                                                                     LPUSH: rps=7740.0 (overall: 16098.7) avg_msec=6.447 (overall: 3.072)                                                                     LPUSH: rps=6956.2 (overall: 14788.1) avg_msec=7.147 (overall: 3.347)                                                                     LPUSH: rps=6652.0 (overall: 13771.6) avg_msec=7.494 (overall: 3.597)                                                                     LPUSH: rps=6484.0 (overall: 12962.2) avg_msec=7.716 (overall: 3.826)                                                                     LPUSH: rps=6204.0 (overall: 12286.7) avg_msec=8.051 (overall: 4.039)                                                                     LPUSH: rps=6024.0 (overall: 11717.6) avg_msec=8.297 (overall: 4.238)                                                                     LPUSH: rps=5724.0 (overall: 11218.3) avg_msec=8.706 (overall: 4.428)                                                                     LPUSH: rps=5592.0 (overall: 10785.6) avg_msec=8.942 (overall: 4.608)                                                                     LPUSH: rps=5024.0 (overall: 10374.2) avg_msec=9.903 (overall: 4.791)                                                                     LPUSH: rps=5160.0 (overall: 10026.7) avg_msec=9.706 (overall: 4.960)                                                                     LPUSH: rps=5255.0 (overall: 9727.4) avg_msec=9.464 (overall: 5.112)                                                                    LPUSH: rps=5156.0 (overall: 9458.6) avg_msec=9.731 (overall: 5.260)                                                                    LPUSH: rps=5184.0 (overall: 9221.2) avg_msec=9.641 (overall: 5.397)                                                                    LPUSH: rps=5104.0 (overall: 9004.6) avg_msec=9.783 (overall: 5.528)                                                                    LPUSH: rps=5036.0 (overall: 8806.3) avg_msec=9.915 (overall: 5.653)                                                                    LPUSH: rps=5044.0 (overall: 8627.2) avg_msec=9.875 (overall: 5.771)                                                                    LPUSH: rps=4605.6 (overall: 8443.8) avg_msec=10.835 (overall: 5.897)                                                                     LPUSH: rps=4460.0 (overall: 8270.6) avg_msec=11.175 (overall: 6.021)                                                                     LPUSH: rps=3724.0 (overall: 8081.3) avg_msec=13.381 (overall: 6.162)                                                                     LPUSH: rps=3812.0 (overall: 7910.6) avg_msec=13.096 (overall: 6.295)                                                                     LPUSH: rps=4164.0 (overall: 7766.6) avg_msec=12.033 (overall: 6.414)                                                                     LPUSH: rps=4072.0 (overall: 7629.8) avg_msec=12.286 (overall: 6.530)                                                                     LPUSH: rps=3880.5 (overall: 7495.4) avg_msec=12.748 (overall: 6.645)                                                                     LPUSH: rps=3772.0 (overall: 7367.1) avg_msec=13.303 (overall: 6.763)                                                                     LPUSH: rps=3552.0 (overall: 7240.0) avg_msec=13.986 (overall: 6.881)                                                                     LPUSH: rps=3668.0 (overall: 7124.8) avg_msec=13.698 (overall: 6.994)                                                                     LPUSH: rps=3624.0 (overall: 7015.5) avg_msec=13.766 (overall: 7.103)                                                                     LPUSH: rps=3609.6 (overall: 6911.9) avg_msec=13.852 (overall: 7.210)                                                                     LPUSH: rps=3760.0 (overall: 6819.3) avg_msec=13.253 (overall: 7.308)                                                                     LPUSH: rps=3852.0 (overall: 6734.6) avg_msec=13.051 (overall: 7.402)                                                                     LPUSH: rps=3588.0 (overall: 6647.2) avg_msec=13.872 (overall: 7.499)                                                                     LPUSH: rps=3460.0 (overall: 6561.1) avg_msec=14.540 (overall: 7.599)                                                                     LPUSH: rps=3564.0 (overall: 6482.3) avg_msec=13.992 (overall: 7.692)                                                                     LPUSH: rps=3560.0 (overall: 6407.4) avg_msec=14.020 (overall: 7.782)                                                                     LPUSH: rps=3342.6 (overall: 6330.5) avg_msec=14.815 (overall: 7.875)                                                                     LPUSH: rps=3240.0 (overall: 6255.2) avg_msec=15.450 (overall: 7.971)                                                                     LPUSH: rps=3380.0 (overall: 6186.8) avg_msec=14.820 (overall: 8.060)                                                                     LPUSH: rps=3316.0 (overall: 6120.0) avg_msec=15.057 (overall: 8.148)                                                                     LPUSH: rps=3300.0 (overall: 6056.0) avg_msec=15.153 (overall: 8.234)                                                                     LPUSH: rps=3152.0 (overall: 5991.5) avg_msec=15.935 (overall: 8.324)                                                                     LPUSH: rps=3168.0 (overall: 5930.1) avg_msec=15.728 (overall: 8.410)                                                                     LPUSH: rps=3087.6 (overall: 5869.4) avg_msec=16.175 (overall: 8.498)                                                                     LPUSH: rps=3116.0 (overall: 5812.1) avg_msec=15.903 (overall: 8.580)                                                                     LPUSH: rps=3008.0 (overall: 5754.9) avg_msec=16.716 (overall: 8.667)                                                                     LPUSH: rps=2888.0 (overall: 5697.6) avg_msec=17.237 (overall: 8.754)                                                                     LPUSH: rps=2900.0 (overall: 5642.8) avg_msec=17.338 (overall: 8.840)                                                                     LPUSH: rps=2884.5 (overall: 5589.6) avg_msec=17.189 (overall: 8.923)                                                                     LPUSH: rps=2964.0 (overall: 5540.1) avg_msec=16.889 (overall: 9.004)                                                                     LPUSH: rps=2936.0 (overall: 5491.9) avg_msec=17.053 (overall: 9.083)                                                                     LPUSH: rps=2968.0 (overall: 5446.0) avg_msec=16.778 (overall: 9.160)                                                                     LPUSH: rps=2940.0 (overall: 5401.3) avg_msec=17.107 (overall: 9.237)                                                                     LPUSH: rps=2784.0 (overall: 5355.4) avg_msec=17.940 (overall: 9.316)                                                                     LPUSH: rps=2620.0 (overall: 5308.2) avg_msec=18.996 (overall: 9.398)                                                                     LPUSH: rps=2788.8 (overall: 5265.4) avg_msec=17.885 (overall: 9.475)                                                                     LPUSH: rps=2700.0 (overall: 5222.7) avg_msec=18.593 (overall: 9.553)                                                                     LPUSH: rps=2660.0 (overall: 5180.7) avg_msec=18.739 (overall: 9.631)                                                                     LPUSH: rps=2788.0 (overall: 5142.1) avg_msec=17.927 (overall: 9.703)                                                                     LPUSH: rps=2852.0 (overall: 5105.8) avg_msec=17.589 (overall: 9.773)                                                                     LPUSH: rps=2629.5 (overall: 5067.0) avg_msec=18.840 (overall: 9.847)                                                                     LPUSH: rps=2660.0 (overall: 5030.0) avg_msec=18.721 (overall: 9.919)                                                                     LPUSH: rps=2380.0 (overall: 4989.8) avg_msec=21.023 (overall: 9.999)                                                                     LPUSH: rps=2736.0 (overall: 4956.2) avg_msec=18.338 (overall: 10.068)                                                                      LPUSH: rps=2656.0 (overall: 4922.4) avg_msec=18.905 (overall: 10.138)                                                                      LPUSH: rps=2593.6 (overall: 4888.5) avg_msec=19.152 (overall: 10.207)                                                                      LPUSH: rps=2632.0 (overall: 4856.3) avg_msec=18.934 (overall: 10.275)                                                                      LPUSH: rps=2640.0 (overall: 4825.1) avg_msec=18.961 (overall: 10.342)                                                                      LPUSH: rps=2568.0 (overall: 4793.8) avg_msec=19.449 (overall: 10.410)                                                                      LPUSH: rps=2500.0 (overall: 4762.4) avg_msec=20.120 (overall: 10.479)                                                                      LPUSH: rps=2492.0 (overall: 4731.7) avg_msec=20.017 (overall: 10.547)                                                                      LPUSH: rps=2486.1 (overall: 4701.7) avg_msec=19.900 (overall: 10.613)                                                                      LPUSH: rps=2504.0 (overall: 4672.8) avg_msec=20.033 (overall: 10.680)                                                                      LPUSH: rps=2432.0 (overall: 4643.7) avg_msec=20.494 (overall: 10.746)                                                                      LPUSH: rps=2460.0 (overall: 4615.7) avg_msec=20.386 (overall: 10.812)                                                                      LPUSH: rps=2452.0 (overall: 4588.4) avg_msec=20.436 (overall: 10.877)                                                                      LPUSH: rps=2458.2 (overall: 4561.6) avg_msec=20.210 (overall: 10.940)                                                                      LPUSH: rps=2504.0 (overall: 4536.2) avg_msec=19.983 (overall: 11.002)                                                                      LPUSH: rps=2368.0 (overall: 4509.8) avg_msec=21.243 (overall: 11.068)                                                                      LPUSH: rps=2396.0 (overall: 4484.4) avg_msec=20.848 (overall: 11.130)                                                                      LPUSH: rps=2356.0 (overall: 4459.0) avg_msec=21.069 (overall: 11.193)                                                                      LPUSH: rps=2384.0 (overall: 4434.7) avg_msec=21.057 (overall: 11.255)                                                                      LPUSH: rps=2350.6 (overall: 4410.3) avg_msec=21.141 (overall: 11.317)                                                                      LPUSH: rps=2352.0 (overall: 4386.7) avg_msec=21.318 (overall: 11.378)                                                                      LPUSH: rps=2304.0 (overall: 4363.0) avg_msec=21.618 (overall: 11.440)                                                                      LPUSH: rps=2316.0 (overall: 4340.1) avg_msec=21.631 (overall: 11.501)                                                                      LPUSH: rps=2252.0 (overall: 4316.9) avg_msec=22.211 (overall: 11.563)                                                                      LPUSH: rps=2314.7 (overall: 4294.8) avg_msec=21.567 (overall: 11.622)                                                                      LPUSH: rps=2300.0 (overall: 4273.1) avg_msec=21.699 (overall: 11.681)                                                                      LPUSH: rps=2308.0 (overall: 4252.0) avg_msec=21.761 (overall: 11.740)                                                                      LPUSH: rps=2292.0 (overall: 4231.2) avg_msec=21.752 (overall: 11.798)                                                                      ====== LPUSH ======
  100000 requests completed in 23.75 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.023 milliseconds (cumulative count 3)
50.000% <= 12.183 milliseconds (cumulative count 50014)
75.000% <= 17.375 milliseconds (cumulative count 75008)
87.500% <= 20.047 milliseconds (cumulative count 87565)
93.750% <= 21.247 milliseconds (cumulative count 93755)
96.875% <= 22.063 milliseconds (cumulative count 96894)
98.438% <= 22.751 milliseconds (cumulative count 98466)
99.219% <= 23.471 milliseconds (cumulative count 99225)
99.609% <= 24.159 milliseconds (cumulative count 99614)
99.805% <= 25.023 milliseconds (cumulative count 99809)
99.902% <= 26.127 milliseconds (cumulative count 99903)
99.951% <= 27.887 milliseconds (cumulative count 99952)
99.976% <= 31.711 milliseconds (cumulative count 99976)
99.988% <= 36.031 milliseconds (cumulative count 99988)
99.994% <= 39.039 milliseconds (cumulative count 99994)
99.997% <= 41.631 milliseconds (cumulative count 99997)
99.998% <= 42.271 milliseconds (cumulative count 99999)
99.999% <= 44.383 milliseconds (cumulative count 100000)
100.000% <= 44.383 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
1.160% <= 0.103 milliseconds (cumulative count 1160)
1.923% <= 0.207 milliseconds (cumulative count 1923)
2.621% <= 0.303 milliseconds (cumulative count 2621)
3.353% <= 0.407 milliseconds (cumulative count 3353)
3.569% <= 0.503 milliseconds (cumulative count 3569)
3.789% <= 0.607 milliseconds (cumulative count 3789)
4.087% <= 0.703 milliseconds (cumulative count 4087)
4.439% <= 0.807 milliseconds (cumulative count 4439)
4.675% <= 0.903 milliseconds (cumulative count 4675)
4.976% <= 1.007 milliseconds (cumulative count 4976)
5.564% <= 1.103 milliseconds (cumulative count 5564)
6.496% <= 1.207 milliseconds (cumulative count 6496)
7.110% <= 1.303 milliseconds (cumulative count 7110)
7.593% <= 1.407 milliseconds (cumulative count 7593)
8.087% <= 1.503 milliseconds (cumulative count 8087)
8.551% <= 1.607 milliseconds (cumulative count 8551)
8.891% <= 1.703 milliseconds (cumulative count 8891)
9.208% <= 1.807 milliseconds (cumulative count 9208)
9.481% <= 1.903 milliseconds (cumulative count 9481)
9.690% <= 2.007 milliseconds (cumulative count 9690)
9.868% <= 2.103 milliseconds (cumulative count 9868)
12.763% <= 3.103 milliseconds (cumulative count 12763)
15.410% <= 4.103 milliseconds (cumulative count 15410)
19.284% <= 5.103 milliseconds (cumulative count 19284)
22.544% <= 6.103 milliseconds (cumulative count 22544)
25.702% <= 7.103 milliseconds (cumulative count 25702)
30.102% <= 8.103 milliseconds (cumulative count 30102)
36.441% <= 9.103 milliseconds (cumulative count 36441)
42.549% <= 10.103 milliseconds (cumulative count 42549)
46.370% <= 11.103 milliseconds (cumulative count 46370)
49.734% <= 12.103 milliseconds (cumulative count 49734)
53.975% <= 13.103 milliseconds (cumulative count 53975)
59.108% <= 14.103 milliseconds (cumulative count 59108)
64.102% <= 15.103 milliseconds (cumulative count 64102)
68.946% <= 16.103 milliseconds (cumulative count 68946)
73.705% <= 17.103 milliseconds (cumulative count 73705)
78.537% <= 18.111 milliseconds (cumulative count 78537)
83.185% <= 19.103 milliseconds (cumulative count 83185)
87.923% <= 20.111 milliseconds (cumulative count 87923)
93.087% <= 21.103 milliseconds (cumulative count 93087)
97.033% <= 22.111 milliseconds (cumulative count 97033)
98.900% <= 23.103 milliseconds (cumulative count 98900)
99.595% <= 24.111 milliseconds (cumulative count 99595)
99.819% <= 25.103 milliseconds (cumulative count 99819)
99.902% <= 26.111 milliseconds (cumulative count 99902)
99.938% <= 27.103 milliseconds (cumulative count 99938)
99.960% <= 28.111 milliseconds (cumulative count 99960)
99.969% <= 29.103 milliseconds (cumulative count 99969)
99.972% <= 30.111 milliseconds (cumulative count 99972)
99.973% <= 31.103 milliseconds (cumulative count 99973)
99.976% <= 32.111 milliseconds (cumulative count 99976)
99.981% <= 33.119 milliseconds (cumulative count 99981)
99.983% <= 34.111 milliseconds (cumulative count 99983)
99.985% <= 35.103 milliseconds (cumulative count 99985)
99.988% <= 36.127 milliseconds (cumulative count 99988)
99.989% <= 37.119 milliseconds (cumulative count 99989)
99.991% <= 38.111 milliseconds (cumulative count 99991)
99.994% <= 39.103 milliseconds (cumulative count 99994)
99.996% <= 40.127 milliseconds (cumulative count 99996)
99.998% <= 42.111 milliseconds (cumulative count 99998)
99.999% <= 43.103 milliseconds (cumulative count 99999)
100.000% <= 45.119 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 4211.41 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
       11.852     0.016    12.183    21.551    23.199    44.383
 RPUSH: rps=2060.0 (overall: 27105.3) avg_msec=1.460 (overall: 1.460)                                                                     RPUSH: rps=94450.2 (overall: 89711.1) avg_msec=0.296 (overall: 0.320)                                                                      RPUSH: rps=93156.0 (overall: 91367.3) avg_msec=0.295 (overall: 0.308)                                                                      RPUSH: rps=89428.0 (overall: 90737.7) avg_msec=0.310 (overall: 0.308)                                                                      RPUSH: rps=97108.0 (overall: 92299.0) avg_msec=0.291 (overall: 0.304)                                                                      ====== RPUSH ======
  100000 requests completed in 1.08 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.039 milliseconds (cumulative count 1)
50.000% <= 0.295 milliseconds (cumulative count 57112)
75.000% <= 0.319 milliseconds (cumulative count 79274)
87.500% <= 0.343 milliseconds (cumulative count 89338)
93.750% <= 0.375 milliseconds (cumulative count 94537)
96.875% <= 0.423 milliseconds (cumulative count 96972)
98.438% <= 0.599 milliseconds (cumulative count 98447)
99.219% <= 1.191 milliseconds (cumulative count 99256)
99.609% <= 1.303 milliseconds (cumulative count 99613)
99.805% <= 1.719 milliseconds (cumulative count 99807)
99.902% <= 2.223 milliseconds (cumulative count 99903)
99.951% <= 2.351 milliseconds (cumulative count 99953)
99.976% <= 2.407 milliseconds (cumulative count 99977)
99.988% <= 2.455 milliseconds (cumulative count 99989)
99.994% <= 2.543 milliseconds (cumulative count 99995)
99.997% <= 2.583 milliseconds (cumulative count 99997)
99.998% <= 2.607 milliseconds (cumulative count 99999)
99.999% <= 2.623 milliseconds (cumulative count 100000)
100.000% <= 2.623 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.014% <= 0.103 milliseconds (cumulative count 14)
4.441% <= 0.207 milliseconds (cumulative count 4441)
66.066% <= 0.303 milliseconds (cumulative count 66066)
96.445% <= 0.407 milliseconds (cumulative count 96445)
98.045% <= 0.503 milliseconds (cumulative count 98045)
98.466% <= 0.607 milliseconds (cumulative count 98466)
98.681% <= 0.703 milliseconds (cumulative count 98681)
98.870% <= 0.807 milliseconds (cumulative count 98870)
98.974% <= 0.903 milliseconds (cumulative count 98974)
99.037% <= 1.007 milliseconds (cumulative count 99037)
99.076% <= 1.103 milliseconds (cumulative count 99076)
99.334% <= 1.207 milliseconds (cumulative count 99334)
99.613% <= 1.303 milliseconds (cumulative count 99613)
99.752% <= 1.407 milliseconds (cumulative count 99752)
99.776% <= 1.503 milliseconds (cumulative count 99776)
99.791% <= 1.607 milliseconds (cumulative count 99791)
99.801% <= 1.703 milliseconds (cumulative count 99801)
99.832% <= 1.807 milliseconds (cumulative count 99832)
99.867% <= 1.903 milliseconds (cumulative count 99867)
99.888% <= 2.007 milliseconds (cumulative count 99888)
99.900% <= 2.103 milliseconds (cumulative count 99900)
100.000% <= 3.103 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 92506.94 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.303     0.032     0.295     0.383     0.935     2.623
 LPOP: rps=70434.3 (overall: 94037.2) avg_msec=0.294 (overall: 0.294)                                                                     LPOP: rps=92976.0 (overall: 93431.5) avg_msec=0.291 (overall: 0.293)                                                                     LPOP: rps=94408.0 (overall: 93786.3) avg_msec=0.288 (overall: 0.291)                                                                     LPOP: rps=95008.0 (overall: 94111.9) avg_msec=0.288 (overall: 0.290)                                                                     ====== LPOP ======
  100000 requests completed in 1.06 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.055 milliseconds (cumulative count 2)
50.000% <= 0.295 milliseconds (cumulative count 57223)
75.000% <= 0.311 milliseconds (cumulative count 75652)
87.500% <= 0.335 milliseconds (cumulative count 90422)
93.750% <= 0.351 milliseconds (cumulative count 94699)
96.875% <= 0.367 milliseconds (cumulative count 97043)
98.438% <= 0.391 milliseconds (cumulative count 98657)
99.219% <= 0.415 milliseconds (cumulative count 99259)
99.609% <= 0.455 milliseconds (cumulative count 99631)
99.805% <= 0.551 milliseconds (cumulative count 99810)
99.902% <= 0.639 milliseconds (cumulative count 99904)
99.951% <= 0.735 milliseconds (cumulative count 99954)
99.976% <= 0.807 milliseconds (cumulative count 99976)
99.988% <= 0.855 milliseconds (cumulative count 99988)
99.994% <= 0.903 milliseconds (cumulative count 99994)
99.997% <= 0.927 milliseconds (cumulative count 99997)
99.998% <= 0.999 milliseconds (cumulative count 99999)
99.999% <= 1.047 milliseconds (cumulative count 100000)
100.000% <= 1.047 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.010% <= 0.103 milliseconds (cumulative count 10)
3.531% <= 0.207 milliseconds (cumulative count 3531)
67.282% <= 0.303 milliseconds (cumulative count 67282)
99.118% <= 0.407 milliseconds (cumulative count 99118)
99.747% <= 0.503 milliseconds (cumulative count 99747)
99.877% <= 0.607 milliseconds (cumulative count 99877)
99.938% <= 0.703 milliseconds (cumulative count 99938)
99.976% <= 0.807 milliseconds (cumulative count 99976)
99.994% <= 0.903 milliseconds (cumulative count 99994)
99.999% <= 1.007 milliseconds (cumulative count 99999)
100.000% <= 1.103 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 94250.71 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.289     0.048     0.295     0.359     0.407     1.047
 RPOP: rps=47621.5 (overall: 94865.1) avg_msec=0.296 (overall: 0.296)                                                                     RPOP: rps=97700.0 (overall: 96750.0) avg_msec=0.281 (overall: 0.286)                                                                     RPOP: rps=104556.0 (overall: 99867.4) avg_msec=0.267 (overall: 0.278)                                                                      RPOP: rps=95892.0 (overall: 98732.9) avg_msec=0.290 (overall: 0.282)                                                                     ====== RPOP ======
  100000 requests completed in 1.02 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.071 milliseconds (cumulative count 1)
50.000% <= 0.287 milliseconds (cumulative count 53986)
75.000% <= 0.311 milliseconds (cumulative count 77784)
87.500% <= 0.335 milliseconds (cumulative count 89964)
93.750% <= 0.351 milliseconds (cumulative count 93881)
96.875% <= 0.383 milliseconds (cumulative count 97443)
98.438% <= 0.407 milliseconds (cumulative count 98536)
99.219% <= 0.455 milliseconds (cumulative count 99265)
99.609% <= 0.535 milliseconds (cumulative count 99623)
99.805% <= 0.647 milliseconds (cumulative count 99809)
99.902% <= 0.783 milliseconds (cumulative count 99905)
99.951% <= 0.895 milliseconds (cumulative count 99955)
99.976% <= 0.975 milliseconds (cumulative count 99979)
99.988% <= 1.015 milliseconds (cumulative count 99991)
99.994% <= 1.039 milliseconds (cumulative count 99994)
99.997% <= 1.103 milliseconds (cumulative count 99997)
99.998% <= 1.175 milliseconds (cumulative count 99999)
99.999% <= 1.199 milliseconds (cumulative count 100000)
100.000% <= 1.199 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.012% <= 0.103 milliseconds (cumulative count 12)
8.090% <= 0.207 milliseconds (cumulative count 8090)
71.007% <= 0.303 milliseconds (cumulative count 71007)
98.536% <= 0.407 milliseconds (cumulative count 98536)
99.528% <= 0.503 milliseconds (cumulative count 99528)
99.768% <= 0.607 milliseconds (cumulative count 99768)
99.854% <= 0.703 milliseconds (cumulative count 99854)
99.919% <= 0.807 milliseconds (cumulative count 99919)
99.959% <= 0.903 milliseconds (cumulative count 99959)
99.986% <= 1.007 milliseconds (cumulative count 99986)
99.997% <= 1.103 milliseconds (cumulative count 99997)
100.000% <= 1.207 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 98231.83 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.282     0.064     0.287     0.359     0.431     1.199

```

## Set & Hash Operations (100000 requests, 50 parallel clients)

```
WARNING: Could not fetch server CONFIG
 SADD: rps=0.0 (overall: 0.0) avg_msec=0.000 (overall: 0.000)                                                             SADD: rps=94628.0 (overall: 94628.0) avg_msec=0.295 (overall: 0.295)                                                                     SADD: rps=91788.0 (overall: 93208.0) avg_msec=0.298 (overall: 0.297)                                                                     SADD: rps=94286.9 (overall: 93568.6) avg_msec=0.290 (overall: 0.294)                                                                     SADD: rps=94404.0 (overall: 93777.2) avg_msec=0.289 (overall: 0.293)                                                                     ====== SADD ======
  100000 requests completed in 1.07 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.079 milliseconds (cumulative count 1)
50.000% <= 0.295 milliseconds (cumulative count 55525)
75.000% <= 0.319 milliseconds (cumulative count 80340)
87.500% <= 0.335 milliseconds (cumulative count 88635)
93.750% <= 0.359 milliseconds (cumulative count 94463)
96.875% <= 0.383 milliseconds (cumulative count 96956)
98.438% <= 0.423 milliseconds (cumulative count 98614)
99.219% <= 0.471 milliseconds (cumulative count 99238)
99.609% <= 0.591 milliseconds (cumulative count 99631)
99.805% <= 0.679 milliseconds (cumulative count 99808)
99.902% <= 0.767 milliseconds (cumulative count 99903)
99.951% <= 0.847 milliseconds (cumulative count 99954)
99.976% <= 0.895 milliseconds (cumulative count 99978)
99.988% <= 0.943 milliseconds (cumulative count 99988)
99.994% <= 0.967 milliseconds (cumulative count 99994)
99.997% <= 0.999 milliseconds (cumulative count 99999)
99.999% <= 1.023 milliseconds (cumulative count 100000)
100.000% <= 1.023 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.004% <= 0.103 milliseconds (cumulative count 4)
3.143% <= 0.207 milliseconds (cumulative count 3143)
65.321% <= 0.303 milliseconds (cumulative count 65321)
98.186% <= 0.407 milliseconds (cumulative count 98186)
99.378% <= 0.503 milliseconds (cumulative count 99378)
99.669% <= 0.607 milliseconds (cumulative count 99669)
99.832% <= 0.703 milliseconds (cumulative count 99832)
99.936% <= 0.807 milliseconds (cumulative count 99936)
99.979% <= 0.903 milliseconds (cumulative count 99979)
99.999% <= 1.007 milliseconds (cumulative count 99999)
100.000% <= 1.103 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 93808.63 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.293     0.072     0.295     0.367     0.447     1.023
 HSET: rps=70748.0 (overall: 96125.0) avg_msec=0.290 (overall: 0.290)                                                                     HSET: rps=99880.0 (overall: 98288.0) avg_msec=0.280 (overall: 0.284)                                                                     HSET: rps=94120.0 (overall: 96764.6) avg_msec=0.299 (overall: 0.290)                                                                     HSET: rps=93768.9 (overall: 95960.4) avg_msec=0.290 (overall: 0.290)                                                                     ====== HSET ======
  100000 requests completed in 1.04 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.055 milliseconds (cumulative count 2)
50.000% <= 0.295 milliseconds (cumulative count 58523)
75.000% <= 0.319 milliseconds (cumulative count 80229)
87.500% <= 0.335 milliseconds (cumulative count 87914)
93.750% <= 0.367 milliseconds (cumulative count 94822)
96.875% <= 0.391 milliseconds (cumulative count 97009)
98.438% <= 0.431 milliseconds (cumulative count 98441)
99.219% <= 0.519 milliseconds (cumulative count 99249)
99.609% <= 0.639 milliseconds (cumulative count 99623)
99.805% <= 0.751 milliseconds (cumulative count 99809)
99.902% <= 0.887 milliseconds (cumulative count 99905)
99.951% <= 0.967 milliseconds (cumulative count 99953)
99.976% <= 1.031 milliseconds (cumulative count 99979)
99.988% <= 1.055 milliseconds (cumulative count 99988)
99.994% <= 1.095 milliseconds (cumulative count 99994)
99.997% <= 1.183 milliseconds (cumulative count 99997)
99.998% <= 1.231 milliseconds (cumulative count 99999)
99.999% <= 1.287 milliseconds (cumulative count 100000)
100.000% <= 1.287 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.048% <= 0.103 milliseconds (cumulative count 48)
5.449% <= 0.207 milliseconds (cumulative count 5449)
67.050% <= 0.303 milliseconds (cumulative count 67050)
97.742% <= 0.407 milliseconds (cumulative count 97742)
99.167% <= 0.503 milliseconds (cumulative count 99167)
99.552% <= 0.607 milliseconds (cumulative count 99552)
99.736% <= 0.703 milliseconds (cumulative count 99736)
99.862% <= 0.807 milliseconds (cumulative count 99862)
99.915% <= 0.903 milliseconds (cumulative count 99915)
99.967% <= 1.007 milliseconds (cumulative count 99967)
99.994% <= 1.103 milliseconds (cumulative count 99994)
99.997% <= 1.207 milliseconds (cumulative count 99997)
100.000% <= 1.303 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 95785.44 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.290     0.048     0.295     0.375     0.487     1.287

```

## Pipelined Operations (100000 requests, 50 parallel clients, 16 pipeline)

```
WARNING: Could not fetch server CONFIG
 SET: rps=0.0 (overall: 0.0) avg_msec=0.000 (overall: 0.000)                                                            ====== SET ======
  100000 requests completed in 0.19 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.095 milliseconds (cumulative count 32)
50.000% <= 1.183 milliseconds (cumulative count 50544)
75.000% <= 1.383 milliseconds (cumulative count 75520)
87.500% <= 1.559 milliseconds (cumulative count 87632)
93.750% <= 1.815 milliseconds (cumulative count 93872)
96.875% <= 2.119 milliseconds (cumulative count 96880)
98.438% <= 2.463 milliseconds (cumulative count 98448)
99.219% <= 2.727 milliseconds (cumulative count 99248)
99.609% <= 3.111 milliseconds (cumulative count 99616)
99.805% <= 3.527 milliseconds (cumulative count 99808)
99.902% <= 3.639 milliseconds (cumulative count 99904)
99.951% <= 3.823 milliseconds (cumulative count 99952)
99.976% <= 4.071 milliseconds (cumulative count 99984)
99.988% <= 5.599 milliseconds (cumulative count 100000)
100.000% <= 5.599 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.032% <= 0.103 milliseconds (cumulative count 32)
0.352% <= 0.207 milliseconds (cumulative count 352)
0.976% <= 0.303 milliseconds (cumulative count 976)
2.480% <= 0.407 milliseconds (cumulative count 2480)
4.576% <= 0.503 milliseconds (cumulative count 4576)
7.664% <= 0.607 milliseconds (cumulative count 7664)
11.792% <= 0.703 milliseconds (cumulative count 11792)
17.440% <= 0.807 milliseconds (cumulative count 17440)
22.864% <= 0.903 milliseconds (cumulative count 22864)
30.496% <= 1.007 milliseconds (cumulative count 30496)
40.448% <= 1.103 milliseconds (cumulative count 40448)
54.080% <= 1.207 milliseconds (cumulative count 54080)
65.872% <= 1.303 milliseconds (cumulative count 65872)
77.824% <= 1.407 milliseconds (cumulative count 77824)
84.688% <= 1.503 milliseconds (cumulative count 84688)
89.600% <= 1.607 milliseconds (cumulative count 89600)
92.032% <= 1.703 milliseconds (cumulative count 92032)
93.744% <= 1.807 milliseconds (cumulative count 93744)
95.056% <= 1.903 milliseconds (cumulative count 95056)
96.096% <= 2.007 milliseconds (cumulative count 96096)
96.784% <= 2.103 milliseconds (cumulative count 96784)
99.584% <= 3.103 milliseconds (cumulative count 99584)
99.984% <= 4.103 milliseconds (cumulative count 99984)
100.000% <= 6.103 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 512820.53 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        1.185     0.088     1.183     1.903     2.663     5.599
 GET: rps=109228.0 (overall: 496490.9) avg_msec=1.398 (overall: 1.398)                                                                      ====== GET ======
  100000 requests completed in 0.19 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.071 milliseconds (cumulative count 16)
50.000% <= 1.343 milliseconds (cumulative count 50224)
75.000% <= 1.495 milliseconds (cumulative count 75344)
87.500% <= 1.663 milliseconds (cumulative count 87792)
93.750% <= 1.919 milliseconds (cumulative count 93856)
96.875% <= 2.167 milliseconds (cumulative count 96928)
98.438% <= 2.567 milliseconds (cumulative count 98448)
99.219% <= 3.063 milliseconds (cumulative count 99248)
99.609% <= 3.527 milliseconds (cumulative count 99616)
99.805% <= 4.127 milliseconds (cumulative count 99808)
99.902% <= 4.647 milliseconds (cumulative count 99904)
99.951% <= 5.047 milliseconds (cumulative count 99952)
99.976% <= 5.479 milliseconds (cumulative count 99984)
99.988% <= 5.743 milliseconds (cumulative count 100000)
100.000% <= 5.743 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.064% <= 0.103 milliseconds (cumulative count 64)
0.560% <= 0.207 milliseconds (cumulative count 560)
1.504% <= 0.303 milliseconds (cumulative count 1504)
2.784% <= 0.407 milliseconds (cumulative count 2784)
4.224% <= 0.503 milliseconds (cumulative count 4224)
5.648% <= 0.607 milliseconds (cumulative count 5648)
7.056% <= 0.703 milliseconds (cumulative count 7056)
8.352% <= 0.807 milliseconds (cumulative count 8352)
10.272% <= 0.903 milliseconds (cumulative count 10272)
13.344% <= 1.007 milliseconds (cumulative count 13344)
17.776% <= 1.103 milliseconds (cumulative count 17776)
26.864% <= 1.207 milliseconds (cumulative count 26864)
42.480% <= 1.303 milliseconds (cumulative count 42480)
62.192% <= 1.407 milliseconds (cumulative count 62192)
76.416% <= 1.503 milliseconds (cumulative count 76416)
85.104% <= 1.607 milliseconds (cumulative count 85104)
89.328% <= 1.703 milliseconds (cumulative count 89328)
91.888% <= 1.807 milliseconds (cumulative count 91888)
93.616% <= 1.903 milliseconds (cumulative count 93616)
95.424% <= 2.007 milliseconds (cumulative count 95424)
96.432% <= 2.103 milliseconds (cumulative count 96432)
99.328% <= 3.103 milliseconds (cumulative count 99328)
99.792% <= 4.103 milliseconds (cumulative count 99792)
99.952% <= 5.103 milliseconds (cumulative count 99952)
100.000% <= 6.103 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 512820.53 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        1.345     0.064     1.343     1.983     2.871     5.743

```

Killing kivo server (PID: 4748)
