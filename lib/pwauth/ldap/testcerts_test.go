package ldap

const rootCAPem = `-----BEGIN CERTIFICATE-----
MIIE1jCCAr4CAQowDQYJKoZIhvcNAQELBQAwMTELMAkGA1UEBhMCVVMxEDAOBgNV
BAoMB1Rlc3RPcmcxEDAOBgNVBAsMB1Rlc3QgQ0EwHhcNMjEwMTI5MTg0MzAyWhcN
NDEwMTI0MTg0MzAyWjAxMQswCQYDVQQGEwJVUzEQMA4GA1UECgwHVGVzdE9yZzEQ
MA4GA1UECwwHVGVzdCBDQTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIB
ANhU96gNK7t3c0EzIGIZkpGtX4l00sP1jNnO3uIeGaVLJZUJ9SJk9Lszqr65oS4Q
NdWPX3uELdv+jBX4716W2xTJfjJADG1qOGl5iFPLB3BTyYOFthU88KOr6u9m22QQ
9aLQEydcQbjwoajC5LKcRpBVb2g1etDSG1vag4OXPmoPcz3Xg2GscLJ3WAnh++ij
W71ZTcbppZa4K7cfFQHlXUOeLAR5biceeAkazcrZyZ69GfnrRjON+RLYbjpkB64S
NAeBXvQnnydPVDkE3t6DJUI+7EOziTI42m06q1XR3XlZRJOHm1SAL5QpKA9CQrec
TuqUk81S+3m/D0rvAHrMAGTyUhpi8pye+IZXeGxGIIsXT1rv0A0xgm0hJLl9BbZL
jtnpFKoxTgYSf3kAhZmn3z7CxUTA7Lma8h+fjIFTKxup05+nM+OMHWrwKVQhzWVk
fnCtk3c03dC85GP50r2hvWhoDD+UT4SKuqxm7E1ung/r4QkQ+JH92s0/w5sUHyhR
IgUsPNqPXHR+/g8LNnm/C3N4TOz26dGScEUAT1YF5e2cEecahyiegXA+fKElWoov
aEav2KQMyQ+fehep8sbuX5i+6jWSO/s4x6ZMaLXpJLzCI4j7ag9fQ+qdjGcilRQY
n8xG9zbE1NXSR8VGvDJXAT7c+SLlsk1EvHyj7hv32FmtAgMBAAEwDQYJKoZIhvcN
AQELBQADggIBABG33C/5XGsO0LRph0a0lQoQoSeYaw9dSc0ThorHCy2ZCntVar8R
A693BQrUujbRC8xE88VzoxTDeMDs9yqfcX7a49EAOlDP5MBss0Ny/iPbCTBHBei6
0nL4pruRct1IdkznNMCJqbPCbkvbF8tDilguKd2L1em3kMODVcO3Ye5yGwLNY6cm
l4C6uqdiCE6Q1qDnUuFYs0oPdUFT8MEmmU1nkFb6eyY2ldJUbk+VJiNBRTXID6ez
vWusUoZq28zOYjAJLfolDODLDUS7YuWoy4NsbXeIyRacqLDvz8dhWDsElMRo7XGO
IUE8neFlW4JopU3I/P82ZeFmEMhDJvIYRYM4lQ0ZLAPORYJwlyRNT6zi0Uyb5C8X
pnFvneXnIIrETPkOOHxN4+Y6CWS0M53nUNzKbcv9fMTEeGi30a6+klds1/yTO06L
PqVFE7K2o0vRHsgYpKe8ZMPG/Y/vnLIU6NSctTz7aWetC3JuGrAFH/XyLofgkV4r
n1qPs587Maf9Ouf5/iVGZUSTxGsTHoVcko/H6feBz1WQ2HMeWwADprv0Mpjglf15
CKK+XDvUZIRFvIRrxEIMEERW0wC3sbN6Drc6bhlth6JQTtYlIxI4iFfC+GPCwKWJ
pg9K87x7CXSKx+wfSt6hG3SQ5fDa1JtlgomzCG5mPEDHBzw5SadE1/Gx
-----END CERTIFICATE-----`

const localhostCertPem = `-----BEGIN CERTIFICATE-----
MIID/DCCAeSgAwIBAgIFAt3gJswwDQYJKoZIhvcNAQELBQAwMTELMAkGA1UEBhMC
VVMxEDAOBgNVBAoMB1Rlc3RPcmcxEDAOBgNVBAsMB1Rlc3QgQ0EwHhcNMjEwMTI5
MTg0MzAyWhcNNDEwMTI0MTg0MzAyWjAUMRIwEAYDVQQDDAlsb2NhbGhvc3QwggEi
MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDAoeVahR73uHgrs0Yj2glLoGPD
mrq6ckQV6nBOw+Ed8BUk6jDn9MBKmRAu/6yuFTjyR04D1a4aCauFfvRco0iSQHvN
Uu6djtzRW05d/txCqC3JnnoblztR1hZhhtwoRs4WvDqsC8VF1P3rhbuIJYpmrPao
y3WQCnGdE+uh4riRgDr3FePn6Rb4aeZjHEXAmo722tHymPHEExl8bGzQKPCDj10W
46vmertgfjNpbm3p8PJnhRVqJ8fGYNXarEfdOR3Rvclt0MIRWdejDXvdl6hG+Q8a
rz7cqjPgH2pOVQSoIFy/b5ezFOXfzCAjum/Lu5heWYnyoDlXOWFgiIoA33RBAgMB
AAGjODA2MAkGA1UdEwQCMAAwFAYDVR0RBA0wC4IJbG9jYWxob3N0MBMGA1UdJQQM
MAoGCCsGAQUFBwMBMA0GCSqGSIb3DQEBCwUAA4ICAQC5UMUkobzs+4uycsZxabVB
jrA+tkKd18z5OM+ZUmWzn8qwf4lWf9C02bLnJ7LFAjmPcTikAL0L+hfeZayKTjIk
5rqqxZDzSsC5BKZmq4rMX5wXoXMzFUksK9rPuXErMgruPvRYBSFNL8rqac1d5WMf
R8hksHHYfe9YF0thZe6kCVD58Mzo0vy10pNbxC0TCAIaDG+OwS8nCbWx3lsoVtq6
a6KAHMdlOaZepJKMBYj9le8fxz6xWaXtjLd/PikyZeKnijDOj8tHGo8pTkMMi9SA
3eQGB5Z/aqIIh0wTIJW0BYjbefUIgp60uae4xeFl4VFD4T6vyZZN7xSmyBd+CZxK
aJk9h8aJLMhvrs8fRSoIRI3mMdOFFK2JSQOvc+VTOekehr/EGgyVkb3jdZg1iNob
l3xUNHByO3DU7QWxGBbCFAR6O5zKqCto5Zq1iEU0nyrlXm6VLtivzLQESW8W47tl
fxdhuwOOZPEHT/mIkEPSMx45FVnaU21q5vxdmayREr6b1xYehQYwjObQgWozB2T/
eRf5ubjv7LZUnz2sdL/FZQRO6Fm2On0HIlpg5I5QVaL5m2/rePazsPpp7NE191m6
iDEbwaHLbzg60J+ighggEAkAilGEZZ0GagtfXWEBkIw3IGgg36+RVcDhzKXV5tNy
/IvhSoSddY8JA/uRROxi0A==
-----END CERTIFICATE-----`

const localhostKeyPem = `-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDAoeVahR73uHgr
s0Yj2glLoGPDmrq6ckQV6nBOw+Ed8BUk6jDn9MBKmRAu/6yuFTjyR04D1a4aCauF
fvRco0iSQHvNUu6djtzRW05d/txCqC3JnnoblztR1hZhhtwoRs4WvDqsC8VF1P3r
hbuIJYpmrPaoy3WQCnGdE+uh4riRgDr3FePn6Rb4aeZjHEXAmo722tHymPHEExl8
bGzQKPCDj10W46vmertgfjNpbm3p8PJnhRVqJ8fGYNXarEfdOR3Rvclt0MIRWdej
DXvdl6hG+Q8arz7cqjPgH2pOVQSoIFy/b5ezFOXfzCAjum/Lu5heWYnyoDlXOWFg
iIoA33RBAgMBAAECggEAeV6H73yoglQMAxy1OKmL6cZolTnMJOUR2O0ZTcdE82Pt
LpEPt1YSQe4msDYPSq+8bYpXsTrUszscgsP2mteWRe+zES8LgOIeZxosSjTl+mmU
T9A2B2RFz84f09rwo7/Y4aI/JV9VMCZ+xgJAogtlJEQeNUPcEqFB7EI82IbM237j
mBINNOJtuAPMIu1VCbEw+kmg2BCtSJOxMheTo2Z+pmmMaGvoSNIm3QV3x4GtjI6I
h/n2iIXUExUGmlz2sWWH92nHr3oai9BXP37NKeGwI0vq4qb130U6UK3MWPpT2ugc
Om9tHeJkhZ6/T8Swv4jd9QLM7d5DOzvX63YxgIJmQQKBgQDnN7h7pDWo4y3Bpu0p
mb5lenKh+cQdncWBYZQv5u9ia9kqIquSW7lvL5JGL+NRuiJi3SP0Sxcs8RN7j5QC
AeL4dCfOU2ViigTtf7ldSCaxiOoa+DpbgYWgf4tCM54dIG9mMV9tHMjF0jZ00zir
nVrWGzeLljJXzdO9ZmQMStGtTQKBgQDVR3OYLsXW6v4rfQMK2SBLCs+6g3LD36tl
04Sa4Xlsrskzz0qeyFDcZwGSOzImSncc2rzLsr2ZUJiOvzsu7YmexoyWb7ErdrLP
NRJsP/213LWp64pJUM6pgsz2/DHzrISaYtc0smwmjOHlQQ6fFJoPYB48pufV6+eE
orAthO14xQKBgQCIEsHOegheWTxvcDa4udNUU2itLJmfOF+o0e1s94LAMGpAouDI
JJUP+zYhekNUsK9V8YEcXyjHeSUXHZtkRwn1YB6hDXFoOYPG5dkILdMfvkzQDHAD
tEkY+JbTIh+WUqVcxge75im+SgVkYX5DeTqhMKlSy9Ta2bYYC+8rUMjvLQKBgQDB
EdiLuDOyVbJXLejWJi38oMHhhuMae90N5ceR6XDhOOy88PcM/CtvCfQ7K0k/roNb
ZIwqHhlSs8oW1vg9iBzf1b8o491PijleKB4QTnFe83ikZKwfqH4cp2LiZvTVMKQt
mjJU6vvKfhh0T0tsKNs59foJT9JpLg+8WwX/fuj2PQKBgQCA37Ksr8xZdgPshb3P
VnFtOKbKBuli+B6XXleigIUx2UdUIq+gjxCaulWjIxygXqDfv81YalV98F+wYVwm
i3Vx30QrwQZx7FcZAQDMdhPgYnnhG6H/oQt5iVrVFeFFNKhTpJEkEpGjDCgASS/U
AnGIeR8ea4lM7li7SZ3qaMgqZA==
-----END PRIVATE KEY-----`


