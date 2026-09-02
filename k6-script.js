import http from "k6/http";
import { check, sleep } from "k6";
import { Rate, Trend } from "k6/metrics";

const BASE_URL = __ENV.BASE_URL || "http://192.168.30.153:8082/api/v1";
const USERNAME = __ENV.API_USERNAME || __ENV.DOKTER_USERNAME || "DRHANDI";
const PASSWORD = __ENV.API_PASSWORD || __ENV.DOKTER_PASSWORD || "1";
const ENDPOINT = (__ENV.ENDPOINT || "all").toLowerCase(); // obat | antrean | pemeriksaan | pemeriksaan_kunjungan | pemeriksaan_pasien | resep | resep_kunjungan | resep_pasien | tindakan_lab | tindakan_lab_detail | tindakan | laboratorium | laboratorium_kunjungan | laboratorium_pasien | all

const errorRate = new Rate("errors");
const obatDuration = new Trend("obat_duration");
const antreanDuration = new Trend("antrean_duration");
const pemeriksaanKunjunganDuration = new Trend("pemeriksaan_kunjungan_duration");
const pemeriksaanPasienDuration = new Trend("pemeriksaan_pasien_duration");
const resepKunjunganDuration = new Trend("resep_kunjungan_duration");
const resepPasienDuration = new Trend("resep_pasien_duration");
const tindakanLabDuration = new Trend("tindakan_lab_duration");
const tindakanLabDetailDuration = new Trend("tindakan_lab_detail_duration");
const laboratoriumKunjunganDuration = new Trend("laboratorium_kunjungan_duration");
const laboratoriumPasienDuration = new Trend("laboratorium_pasien_duration");

export const options = {
    scenarios: {
        operasional: {
            executor: "ramping-vus",
            exec: "mixedWorkload",
            startVUs: 0,
            stages: [
                { duration: "5s", target: 5 },
                { duration: "20s", target: 15 },
                { duration: "5s", target: 0 },
            ],
            tags: { skenario: "operasional" },
        },
        load: {
            executor: "ramping-vus",
            exec: "mixedWorkload",
            startVUs: 0,
            startTime: "35s",
            stages: [
                { duration: "5s", target: 30 },
                { duration: "20s", target: 50 },
                { duration: "5s", target: 0 },
            ],
            tags: { skenario: "load" },
        },
        stress: {
            executor: "ramping-vus",
            exec: "mixedWorkload",
            startVUs: 0,
            startTime: "70s",
            stages: [
                { duration: "5s", target: 80 },
                { duration: "20s", target: 100 },
                { duration: "5s", target: 0 },
            ],
            tags: { skenario: "stress" },
        },
    },
    thresholds: {
        http_req_failed: ["rate<0.01"],
        "http_req_duration{skenario:operasional}": ["p(95)<200"],
        "http_req_duration{skenario:load}":        ["p(95)<500"],
        "http_req_duration{skenario:stress}":      ["p(95)<1000"],
        "errors{skenario:operasional}": ["rate<0.01"],
        "errors{skenario:load}":        ["rate<0.05"],
        "errors{skenario:stress}":      ["rate<0.10"],
    },
};

const statusLanjutList = ["Semua", "Ralan", "Ranap"];
const statusPemeriksaanList = ["Belum", "Sudah"];
const penjaminList = ["BPJ", "UMU"];
const obatKeywords = ["PARA", "AMOX", "INJ", "TAB", "SYR", "OMEP", "CETIR", "DEXA"];
const antreanKeywords = ["Ahmad", "Siti", "Budi", "Dewi", "Rina"];
const kategoriLabList = ["pk", "pa", "mb"];
const labKeywords = ["Darah", "Urin", "Glukosa", "Kolesterol", "Kultur", "Biopsi", "SGOT", "SGPT", "Ureum", "Kreatinin"];
const today = "2026-04-22,2026-04-22";

export function setup() {
    const loginPayload = JSON.stringify({
        username: USERNAME,
        password: PASSWORD,
    });

    const loginParams = {
        headers: { "Content-Type": "application/json" },
    };

    const loginRes = http.post(`${BASE_URL}/auth/login`, loginPayload, loginParams);

    if (loginRes.status !== 200) {
        throw new Error(`Login gagal dengan status ${loginRes.status}: ${loginRes.body}`);
    }

    const token = loginRes.json("data.token");
    const namaDokter = loginRes.json("data.nama_dokter") || USERNAME;

    if (!token) {
        throw new Error(`Token tidak ditemukan dalam response: ${loginRes.body}`);
    }

    console.log(`[SETUP] Login berhasil untuk: ${namaDokter} (${USERNAME})`);

    const authHeaders = {
        headers: {
            "Accept": "application/json",
            "Authorization": `Bearer ${token}`,
        },
    };

    // Fetch antrean untuk mendapatkan id_kunjungan & id_pasien aktif
    const antreanRes = http.get(`${BASE_URL}/rawat-jalan/antrean?tanggal=${today}&page=1&limit=20`, authHeaders);
    let kunjungans = [];
    let pasiens = [];

    if (antreanRes.status === 200) {
        const items = antreanRes.json("data") || [];
        kunjungans = items.map((i) => i.id).filter(Boolean);
        pasiens = items.map((i) => i.id_pasien).filter(Boolean);
    }

    // Fetch sample tindakan lab per kategori untuk detail endpoint test
    const tindakanLabSamples = {};
    for (const kat of kategoriLabList) {
        const labRes = http.get(`${BASE_URL}/tindakan/lab/${kat}?page=1&limit=10`, authHeaders);
        if (labRes.status === 200) {
            const items = labRes.json("data") || [];
            tindakanLabSamples[kat] = items.map((i) => i.id).filter(Boolean);
        } else {
            tindakanLabSamples[kat] = [];
        }
    }

    console.log(`[SETUP] Berhasil mengambil ${kunjungans.length} sampel ID kunjungan & ${pasiens.length} ID pasien.`);
    console.log(`[SETUP] Sampel Tindakan Lab: PK=${(tindakanLabSamples.pk || []).length}, PA=${(tindakanLabSamples.pa || []).length}, MB=${(tindakanLabSamples.mb || []).length}`);
    console.log(`[SETUP] Target Endpoint: ${ENDPOINT}`);

    return {
        token: token,
        kunjungans: kunjungans,
        pasiens: pasiens,
        tindakanLabSamples: tindakanLabSamples,
    };
}

function requestObat(params) {
    const scenario = Math.random();
    let url = "";

    if (scenario < 0.4) {
        const keyword = obatKeywords[Math.floor(Math.random() * obatKeywords.length)];
        url = `${BASE_URL}/obat?keyword=${keyword}&page=1&limit=20`;
    } else if (scenario < 0.7) {
        url = `${BASE_URL}/obat?depo=DPRJ&page=1&limit=20`;
    } else {
        const page = Math.floor(Math.random() * 5) + 1;
        url = `${BASE_URL}/obat?page=${page}&limit=20`;
    }

    const res = http.get(url, params);
    obatDuration.add(res.timings.duration);
    return res;
}

function requestAntrean(params) {
    const scenario = Math.random();
    let url = "";

    if (scenario < 0.3) {
        url = `${BASE_URL}/rawat-jalan/antrean?tanggal=${today}&page=1&limit=20`;
    } else if (scenario < 0.5) {
        const penjamin = penjaminList[Math.floor(Math.random() * penjaminList.length)];
        url = `${BASE_URL}/rawat-jalan/antrean?tanggal=${today}&penjamin=${penjamin}&page=1&limit=20`;
    } else if (scenario < 0.7) {
        const status = statusPemeriksaanList[Math.floor(Math.random() * statusPemeriksaanList.length)];
        url = `${BASE_URL}/rawat-jalan/antrean?tanggal=${today}&status_pemeriksaan=${status}&page=1&limit=20`;
    } else if (scenario < 0.85) {
        const keyword = antreanKeywords[Math.floor(Math.random() * antreanKeywords.length)];
        url = `${BASE_URL}/rawat-jalan/antrean?tanggal=${today}&keyword=${keyword}&page=1&limit=20`;
    } else {
        const page = Math.floor(Math.random() * 4) + 2;
        url = `${BASE_URL}/rawat-jalan/antrean?tanggal=${today}&page=${page}&limit=20`;
    }

    const res = http.get(url, params);
    antreanDuration.add(res.timings.duration);
    return res;
}

function requestPemeriksaanKunjungan(params, data) {
    if (!data.kunjungans || data.kunjungans.length === 0) {
        return requestAntrean(params);
    }

    const idKunjungan = data.kunjungans[Math.floor(Math.random() * data.kunjungans.length)];
    const statusLanjut = statusLanjutList[Math.floor(Math.random() * statusLanjutList.length)];

    const scenario = Math.random();
    let url = `${BASE_URL}/pemeriksaan/${idKunjungan}/${statusLanjut}`;

    if (scenario < 0.5) {
        url += `?page=1&limit=20`;
    } else if (scenario < 0.8) {
        url += `?tanggal=${today}&page=1&limit=20`;
    } else {
        const page = Math.floor(Math.random() * 3) + 1;
        url += `?page=${page}&limit=10`;
    }

    const res = http.get(url, params);
    pemeriksaanKunjunganDuration.add(res.timings.duration);
    return res;
}

function requestPemeriksaanPasien(params, data) {
    if (!data.pasiens || data.pasiens.length === 0) {
        return requestAntrean(params);
    }

    const idPasien = data.pasiens[Math.floor(Math.random() * data.pasiens.length)];
    const statusLanjut = statusLanjutList[Math.floor(Math.random() * statusLanjutList.length)];

    const scenario = Math.random();
    let url = `${BASE_URL}/pemeriksaan/pasien/${idPasien}/${statusLanjut}`;

    if (scenario < 0.5) {
        url += `?page=1&limit=20`;
    } else if (scenario < 0.8) {
        url += `?tanggal=2025-01-01,2026-12-31&page=1&limit=20`;
    } else {
        const page = Math.floor(Math.random() * 4) + 1;
        url += `?page=${page}&limit=20`;
    }

    const res = http.get(url, params);
    pemeriksaanPasienDuration.add(res.timings.duration);
    return res;
}

function requestResepKunjungan(params, data) {
    if (!data.kunjungans || data.kunjungans.length === 0) {
        return requestAntrean(params);
    }

    const idKunjungan = data.kunjungans[Math.floor(Math.random() * data.kunjungans.length)];
    const statusLanjut = statusLanjutList[Math.floor(Math.random() * statusLanjutList.length)];

    const scenario = Math.random();
    let url = `${BASE_URL}/resep/${idKunjungan}/${statusLanjut}`;

    if (scenario < 0.5) {
        url += `?page=1&limit=20`;
    } else if (scenario < 0.8) {
        url += `?tanggal=${today}&page=1&limit=20`;
    } else {
        const page = Math.floor(Math.random() * 3) + 1;
        url += `?page=${page}&limit=10`;
    }

    const res = http.get(url, params);
    resepKunjunganDuration.add(res.timings.duration);
    return res;
}

function requestResepPasien(params, data) {
    if (!data.pasiens || data.pasiens.length === 0) {
        return requestAntrean(params);
    }

    const idPasien = data.pasiens[Math.floor(Math.random() * data.pasiens.length)];
    const statusLanjut = statusLanjutList[Math.floor(Math.random() * statusLanjutList.length)];

    const scenario = Math.random();
    let url = `${BASE_URL}/resep/pasien/${idPasien}/${statusLanjut}`;

    if (scenario < 0.5) {
        url += `?page=1&limit=20`;
    } else if (scenario < 0.8) {
        url += `?tanggal=2025-01-01,2026-12-31&page=1&limit=20`;
    } else {
        const page = Math.floor(Math.random() * 4) + 1;
        url += `?page=${page}&limit=20`;
    }

    const res = http.get(url, params);
    resepPasienDuration.add(res.timings.duration);
    return res;
}

function requestTindakanLab(params) {
    const kategori = kategoriLabList[Math.floor(Math.random() * kategoriLabList.length)];
    const scenario = Math.random();
    let url = `${BASE_URL}/tindakan/lab/${kategori}`;

    if (scenario < 0.5) {
        const keyword = labKeywords[Math.floor(Math.random() * labKeywords.length)];
        url += `?keyword=${keyword}&page=1&limit=20`;
    } else if (scenario < 0.8) {
        const page = Math.floor(Math.random() * 3) + 1;
        url += `?page=${page}&limit=20`;
    } else {
        url += `?page=1&limit=50`;
    }

    const res = http.get(url, params);
    tindakanLabDuration.add(res.timings.duration);
    return res;
}

function requestTindakanLabDetail(params, data) {
    const samples = data.tindakanLabSamples || {};
    const availableCategories = Object.keys(samples).filter((k) => samples[k] && samples[k].length > 0);

    if (availableCategories.length === 0) {
        return requestTindakanLab(params);
    }

    const kat = availableCategories[Math.floor(Math.random() * availableCategories.length)];
    const idList = samples[kat];
    const idTindakan = idList[Math.floor(Math.random() * idList.length)];

    const url = `${BASE_URL}/tindakan/lab/${kat}/${idTindakan}`;
    const res = http.get(url, params);
    tindakanLabDetailDuration.add(res.timings.duration);
    return res;
}

function requestLaboratoriumKunjungan(params, data) {
    if (!data.kunjungans || data.kunjungans.length === 0) {
        return requestTindakanLab(params);
    }

    const idKunjungan = data.kunjungans[Math.floor(Math.random() * data.kunjungans.length)];
    const kategori = Math.random() < 0.6 ? "pk" : "pa";
    const statusLanjut = statusLanjutList[Math.floor(Math.random() * statusLanjutList.length)];

    const url = `${BASE_URL}/laboratorium/${kategori}/${idKunjungan}/${statusLanjut}?page=1&limit=5`;
    const res = http.get(url, params);
    laboratoriumKunjunganDuration.add(res.timings.duration);
    return res;
}

function requestLaboratoriumPasien(params, data) {
    if (!data.pasiens || data.pasiens.length === 0) {
        return requestTindakanLab(params);
    }

    const idPasien = data.pasiens[Math.floor(Math.random() * data.pasiens.length)];
    const kategori = Math.random() < 0.6 ? "pk" : "pa";
    const statusLanjut = statusLanjutList[Math.floor(Math.random() * statusLanjutList.length)];

    const url = `${BASE_URL}/laboratorium/${kategori}/pasien/${idPasien}/${statusLanjut}?page=1&limit=5`;
    const res = http.get(url, params);
    laboratoriumPasienDuration.add(res.timings.duration);
    return res;
}

export default function (data) {
    mixedWorkload(data);
}

export function mixedWorkload(data) {
    const params = {
        headers: {
            "Accept": "application/json",
            "Authorization": `Bearer ${data.token}`,
        },
    };

    let res;
    if (ENDPOINT === "obat") {
        res = requestObat(params);
    } else if (ENDPOINT === "antrean") {
        res = requestAntrean(params);
    } else if (ENDPOINT === "pemeriksaan_kunjungan") {
        res = requestPemeriksaanKunjungan(params, data);
    } else if (ENDPOINT === "pemeriksaan_pasien") {
        res = requestPemeriksaanPasien(params, data);
    } else if (ENDPOINT === "pemeriksaan") {
        res = Math.random() < 0.5 ? requestPemeriksaanKunjungan(params, data) : requestPemeriksaanPasien(params, data);
    } else if (ENDPOINT === "resep_kunjungan") {
        res = requestResepKunjungan(params, data);
    } else if (ENDPOINT === "resep_pasien") {
        res = requestResepPasien(params, data);
    } else if (ENDPOINT === "resep") {
        res = Math.random() < 0.5 ? requestResepKunjungan(params, data) : requestResepPasien(params, data);
    } else if (ENDPOINT === "tindakan_lab") {
        res = requestTindakanLab(params);
    } else if (ENDPOINT === "tindakan_lab_detail") {
        res = requestTindakanLabDetail(params, data);
    } else if (ENDPOINT === "tindakan") {
        res = Math.random() < 0.6 ? requestTindakanLab(params) : requestTindakanLabDetail(params, data);
    } else if (ENDPOINT === "laboratorium_kunjungan") {
        res = requestLaboratoriumKunjungan(params, data);
    } else if (ENDPOINT === "laboratorium_pasien") {
        res = requestLaboratoriumPasien(params, data);
    } else if (ENDPOINT === "laboratorium") {
        res = Math.random() < 0.5 ? requestLaboratoriumKunjungan(params, data) : requestLaboratoriumPasien(params, data);
    } else {
        // Mode 'all': bagi beban ke seluruh modul
        const rand = Math.random();
        if (rand < 0.15) {
            res = requestAntrean(params);
        } else if (rand < 0.30) {
            res = requestObat(params);
        } else if (rand < 0.45) {
            res = requestPemeriksaanKunjungan(params, data);
        } else if (rand < 0.60) {
            res = requestPemeriksaanPasien(params, data);
        } else if (rand < 0.70) {
            res = requestResepKunjungan(params, data);
        } else if (rand < 0.80) {
            res = requestResepPasien(params, data);
        } else if (rand < 0.88) {
            res = requestTindakanLab(params);
        } else if (rand < 0.94) {
            res = requestLaboratoriumKunjungan(params, data);
        } else {
            res = requestLaboratoriumPasien(params, data);
        }
    }

    const isSuccess = check(res, {
        "status is 200": (r) => r.status === 200,
        "response time acceptable": (r) => r.timings.duration < 1000,
        "content type is JSON": (r) =>
            r.headers["Content-Type"] && r.headers["Content-Type"].includes("application/json"),
        "success or status OK": (r) => {
            try {
                const body = r.json();
                return body && (body.code === 200 || body.status === "OK" || body.success === true);
            } catch (_) {
                return false;
            }
        },
    });

    if (!isSuccess && res.status !== 200) {
        console.warn(`[FAIL] HTTP ${res.status} on ${res.url}: ${res.body}`);
    }

    errorRate.add(!isSuccess);
    sleep(0.5);
}

export function teardown() {
    console.log("[TEARDOWN] Performance test selesai.");
}

// Cara menjalankan pengujian k6:
// # Test seluruh endpoint riwayat laboratorium (kunjungan + pasien):
// k6 run -e ENDPOINT=laboratorium k6-script.js
//
// # Test hanya riwayat hasil lab kunjungan:
// k6 run -e ENDPOINT=laboratorium_kunjungan k6-script.js
//
// # Test hanya riwayat hasil lab pasien by RM:
// k6 run -e ENDPOINT=laboratorium_pasien k6-script.js
//
// # Test master tindakan lab:
// k6 run -e ENDPOINT=tindakan k6-script.js
//
// # Test seluruh modul backend ERM Dokter:
// k6 run k6-script.js
