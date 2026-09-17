import http from "k6/http";
import { check, sleep } from "k6";
import { Rate, Trend } from "k6/metrics";

const BASE_URL = __ENV.BASE_URL || "http://localhost:8082/api/v1";
const USERNAME = __ENV.API_USERNAME || __ENV.DOKTER_USERNAME || "DRHANDI";
const PASSWORD = __ENV.API_PASSWORD || __ENV.DOKTER_PASSWORD || "1";
const SERVICE_API_KEY = __ENV.SERVICE_API_KEY || "5c603a34-8254-4e4e-9a25-91fe5a26ff65";
const AUTH_MODE = (__ENV.AUTH_MODE || "jwt").toLowerCase(); // jwt | api_key
const ENDPOINT = (__ENV.ENDPOINT || "all").toLowerCase(); // obat | antrean | antrean_compare | pemeriksaan | pemeriksaan_kunjungan | pemeriksaan_pasien | resep | resep_kunjungan | resep_pasien | tindakan_lab | tindakan_lab_detail | tindakan | laboratorium | laboratorium_kunjungan | laboratorium_pasien | permintaan_lab | permintaan_lab_kunjungan | permintaan_lab_pasien | permintaan_lab_detail | rawat_inap | riwayat_pasien | rawatinap_riwayat | icd10 | icd9 | icd | diagnosa | diagnosa_kunjungan | diagnosa_pasien | all
const PROFILE = (__ENV.PROFILE || "standard").toLowerCase(); // standard | quick | smoke

const errorRate = new Rate("errors");
const obatDuration = new Trend("obat_duration");
const icd10Duration = new Trend("icd10_duration");
const icd9Duration = new Trend("icd9_duration");
const antreanDuration = new Trend("antrean_duration");
const antreanDokterDuration = new Trend("antrean_dokter_duration", true);
const antreanApiKeyDuration = new Trend("antrean_api_key_duration", true);
const antreanDokterFail = new Rate("antrean_dokter_fail");
const antreanApiKeyFail = new Rate("antrean_api_key_fail");

const pemeriksaanKunjunganDuration = new Trend("pemeriksaan_kunjungan_duration");
const pemeriksaanPasienDuration = new Trend("pemeriksaan_pasien_duration");
const resepKunjunganDuration = new Trend("resep_kunjungan_duration");
const resepPasienDuration = new Trend("resep_pasien_duration");
const tindakanLabDuration = new Trend("tindakan_lab_duration");
const tindakanLabDetailDuration = new Trend("tindakan_lab_detail_duration");
const laboratoriumKunjunganDuration = new Trend("laboratorium_kunjungan_duration");
const laboratoriumPasienDuration = new Trend("laboratorium_pasien_duration");
const permintaanLabKunjunganDuration = new Trend("permintaan_lab_kunjungan_duration");
const permintaanLabPasienDuration = new Trend("permintaan_lab_pasien_duration");
const permintaanLabDetailDuration = new Trend("permintaan_lab_detail_duration");
const rawatInapDuration = new Trend("rawat_inap_duration");
const riwayatPasienDuration = new Trend("riwayat_pasien_duration");
const diagnosaKunjunganDuration = new Trend("diagnosa_kunjungan_duration");
const diagnosaPasienDuration = new Trend("diagnosa_pasien_duration");

let scenariosConfig;
let thresholdsConfig;

if (PROFILE === "quick" || PROFILE === "smoke") {
    scenariosConfig = {
        quick_test: {
            executor: "constant-vus",
            exec: (ENDPOINT === "antrean_compare") ? "antreanDokterWorkload" : "mixedWorkload",
            vus: 10,
            duration: "10s",
            tags: { skenario: "quick" },
        },
    };
    thresholdsConfig = {
        http_req_failed: ["rate<0.01"],
        http_req_duration: ["p(95)<150"],
    };
} else if (ENDPOINT === "antrean_compare") {
    scenariosConfig = {
        skenario_dokter_jwt: {
            executor: "ramping-vus",
            exec: "antreanDokterWorkload",
            startVUs: 0,
            stages: [
                { duration: "3s", target: 10 },
                { duration: "10s", target: 20 },
                { duration: "2s", target: 0 },
            ],
            tags: { skenario: "dokter_jwt" },
        },
        skenario_api_key: {
            executor: "ramping-vus",
            exec: "antreanApiKeyWorkload",
            startVUs: 0,
            startTime: "16s",
            stages: [
                { duration: "3s", target: 10 },
                { duration: "10s", target: 20 },
                { duration: "2s", target: 0 },
            ],
            tags: { skenario: "api_key" },
        },
    };
    thresholdsConfig = {
        http_req_failed: ["rate<0.01"],
        "antrean_dokter_duration": ["p(95)<300"],
        "antrean_api_key_duration": ["p(95)<300"],
        "antrean_dokter_fail": ["rate<0.01"],
        "antrean_api_key_fail": ["rate<0.01"],
    };
} else {
    scenariosConfig = {
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
    };
    thresholdsConfig = {
        http_req_failed: ["rate<0.01"],
        "http_req_duration{skenario:operasional}": ["p(95)<200"],
        "http_req_duration{skenario:load}":        ["p(95)<500"],
        "http_req_duration{skenario:stress}":      ["p(95)<1000"],
        "errors{skenario:operasional}": ["rate<0.01"],
        "errors{skenario:load}":        ["rate<0.05"],
        "errors{skenario:stress}":      ["rate<0.10"],
    };
}

export const options = {
    scenarios: scenariosConfig,
    thresholds: thresholdsConfig,
};

const statusLanjutList = ["Semua", "Ralan", "Ranap"];
const statusPemeriksaanList = ["Belum", "Sudah"];
const penjaminList = ["BPJ", "UMU"];
const obatKeywords = ["PARA", "AMOX", "INJ", "TAB", "SYR", "OMEP", "CETIR", "DEXA"];
const antreanKeywords = ["Ahmad", "Siti", "Budi", "Dewi", "Rina"];
const kategoriLabList = ["pk", "pa", "mb"];
const labKeywords = ["Darah", "Urin", "Glukosa", "Kolesterol", "Kultur", "Biopsi", "SGOT", "SGPT", "Ureum", "Kreatinin"];
const icd10Keywords = [
    "stroke", "infark", "febris", "diabetes", "hipertensi", "gastritis", "dyspepsia", 
    "diare", "pneumonia", "asthma", "cephal", "anemia", "I63", "E11", "I10", "A09", "K29", "J18"
];
const icd9Keywords = [
    "tomography", "ultrasound", "radiography", "injection", "infusion", "endoscopy", 
    "transfusion", "dialysis", "biopsy", "87.0", "88.7", "99.0", "96.0", "99.2", "45.1"
];
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
    if (ENDPOINT === "all" || ENDPOINT.includes("lab") || ENDPOINT.includes("tindakan")) {
        for (const kat of kategoriLabList) {
            const labRes = http.get(`${BASE_URL}/tindakan/lab/${kat}?page=1&limit=10`, authHeaders);
            if (labRes.status === 200) {
                const items = labRes.json("data") || [];
                tindakanLabSamples[kat] = items.map((i) => i.id).filter(Boolean);
            } else {
                tindakanLabSamples[kat] = [];
            }
        }
    }

    // Fetch sample permintaan lab PK jika ada
    const permintaanLabSamples = [];
    if (ENDPOINT === "all" || ENDPOINT.includes("lab") || ENDPOINT.includes("permintaan")) {
        for (const idKunj of kunjungans.slice(0, 10)) {
            const pRes = http.get(`${BASE_URL}/laboratorium/pk/permintaan/${idKunj}/Semua`, authHeaders);
            if (pRes.status === 200) {
                const items = pRes.json("data") || [];
                for (const item of items) {
                    if (item.id) {
                        permintaanLabSamples.push({ id_kunjungan: idKunj, id_permintaan: item.id });
                    }
                }
            }
        }
    }

    // Fetch sample pasien dari rawat inap
    const ranapRes = http.get(`${BASE_URL}/rawat-inap/pasien?status=mrs&page=1&limit=20`, authHeaders);
    if (ranapRes.status === 200) {
        const items = ranapRes.json("data") || [];
        for (const item of items) {
            if (item.id_pasien && !pasiens.includes(item.id_pasien)) {
                pasiens.push(item.id_pasien);
            }
        }
    }

    // Fetch master bangsal untuk variasi filter rawat inap
    let bangsalList = [];
    const bangsalRes = http.get(`${BASE_URL}/master/bangsal`, authHeaders);
    if (bangsalRes.status === 200) {
        const items = bangsalRes.json("data") || [];
        bangsalList = items.map((b) => b.kode).filter(Boolean);
    }

    // Fetch master kelas kamar untuk variasi filter rawat inap
    let kelasList = [];
    const kelasRes = http.get(`${BASE_URL}/master/kelas`, authHeaders);
    if (kelasRes.status === 200) {
        const items = kelasRes.json("data") || [];
        kelasList = items.map((k) => k.kode).filter(Boolean);
    }

    console.log(`[SETUP] Berhasil mengambil ${kunjungans.length} sampel ID kunjungan & ${pasiens.length} ID pasien.`);
    console.log(`[SETUP] Master Bangsal: ${bangsalList.length}, Kelas Kamar: ${kelasList.length}`);
    console.log(`[SETUP] Sampel Tindakan Lab: PK=${(tindakanLabSamples.pk || []).length}, PA=${(tindakanLabSamples.pa || []).length}, MB=${(tindakanLabSamples.mb || []).length}`);
    console.log(`[SETUP] Sampel Permintaan Lab PK: ${permintaanLabSamples.length}`);
    console.log(`[SETUP] Target Endpoint: ${ENDPOINT}`);

    return {
        token: token,
        apiKey: SERVICE_API_KEY,
        kunjungans: kunjungans,
        pasiens: pasiens,
        bangsalList: bangsalList,
        kelasList: kelasList,
        tindakanLabSamples: tindakanLabSamples,
        permintaanLabSamples: permintaanLabSamples,
    };
}

function requestICD10(params) {
    const scenario = Math.random();
    let url = "";

    if (scenario < 0.65) {
        const keyword = icd10Keywords[Math.floor(Math.random() * icd10Keywords.length)];
        url = `${BASE_URL}/master/icd10?keyword=${encodeURIComponent(keyword)}&page=1&limit=20`;
    } else if (scenario < 0.85) {
        const page = Math.floor(Math.random() * 5) + 1;
        url = `${BASE_URL}/master/icd10?page=${page}&limit=20`;
    } else {
        url = `${BASE_URL}/master/icd10?page=1&limit=50`;
    }

    const res = http.get(url, params);
    icd10Duration.add(res.timings.duration);
    return res;
}

function requestICD9(params) {
    const scenario = Math.random();
    let url = "";

    if (scenario < 0.65) {
        const keyword = icd9Keywords[Math.floor(Math.random() * icd9Keywords.length)];
        url = `${BASE_URL}/master/icd9?keyword=${encodeURIComponent(keyword)}&page=1&limit=20`;
    } else if (scenario < 0.85) {
        const page = Math.floor(Math.random() * 5) + 1;
        url = `${BASE_URL}/master/icd9?page=${page}&limit=20`;
    } else {
        url = `${BASE_URL}/master/icd9?page=1&limit=50`;
    }

    const res = http.get(url, params);
    icd9Duration.add(res.timings.duration);
    return res;
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

function requestPermintaanLabKunjungan(params, data) {
    if (!data.kunjungans || data.kunjungans.length === 0) {
        return requestTindakanLab(params);
    }

    const idKunjungan = data.kunjungans[Math.floor(Math.random() * data.kunjungans.length)];
    const statusLanjut = statusLanjutList[Math.floor(Math.random() * statusLanjutList.length)];
    const url = `${BASE_URL}/laboratorium/pk/permintaan/${idKunjungan}/${statusLanjut}`;
    const res = http.get(url, params);
    permintaanLabKunjunganDuration.add(res.timings.duration);
    return res;
}

function requestPermintaanLabPasien(params, data) {
    if (!data.pasiens || data.pasiens.length === 0) {
        return requestTindakanLab(params);
    }

    const idPasien = data.pasiens[Math.floor(Math.random() * data.pasiens.length)];
    const statusLanjut = statusLanjutList[Math.floor(Math.random() * statusLanjutList.length)];
    const scenario = Math.random();
    let url = `${BASE_URL}/laboratorium/pk/permintaan/pasien/${idPasien}/${statusLanjut}`;

    if (scenario < 0.5) {
        url += `?page=1&limit=10`;
    } else {
        url += `?tanggal=2025-01-01,2026-12-31&page=1&limit=10`;
    }

    const res = http.get(url, params);
    permintaanLabPasienDuration.add(res.timings.duration);
    return res;
}

function requestPermintaanLabDetail(params, data) {
    const samples = data.permintaanLabSamples || [];
    if (samples.length === 0) {
        return requestPermintaanLabKunjungan(params, data);
    }

    const sample = samples[Math.floor(Math.random() * samples.length)];
    const url = `${BASE_URL}/laboratorium/pk/permintaan/${sample.id_kunjungan}/Semua/${sample.id_permintaan}`;
    const res = http.get(url, params);
    permintaanLabDetailDuration.add(res.timings.duration);
    return res;
}

function requestRawatInap(params, data) {
    const scenario = Math.random();
    let url = "";

    if (scenario < 0.35) {
        url = `${BASE_URL}/rawat-inap/pasien?status=mrs&page=1&limit=20`;
    } else if (scenario < 0.50 && data && data.bangsalList && data.bangsalList.length > 0) {
        const kdBangsal = data.bangsalList[Math.floor(Math.random() * data.bangsalList.length)];
        url = `${BASE_URL}/rawat-inap/pasien?status=mrs&kode_bangsal=${encodeURIComponent(kdBangsal)}&page=1&limit=20`;
    } else if (scenario < 0.65 && data && data.kelasList && data.kelasList.length > 0) {
        const kelas = data.kelasList[Math.floor(Math.random() * data.kelasList.length)];
        url = `${BASE_URL}/rawat-inap/pasien?status=mrs&kelas_kamar=${encodeURIComponent(kelas)}&page=1&limit=20`;
    } else if (scenario < 0.75) {
        const penjamin = penjaminList[Math.floor(Math.random() * penjaminList.length)];
        url = `${BASE_URL}/rawat-inap/pasien?status=mrs&kode_penjamin=${penjamin}&page=1&limit=20`;
    } else if (scenario < 0.85) {
        url = `${BASE_URL}/rawat-inap/pasien?status=krs&tanggal=${today}&page=1&limit=20`;
    } else if (scenario < 0.93) {
        const kw = antreanKeywords[Math.floor(Math.random() * antreanKeywords.length)];
        url = `${BASE_URL}/rawat-inap/pasien?keyword=${encodeURIComponent(kw)}&page=1&limit=20`;
    } else {
        url = `${BASE_URL}/rawat-inap/pasien?status=mrs&page=2&limit=20`;
    }

    const res = http.get(url, params);
    rawatInapDuration.add(res.timings.duration);
    return res;
}

function requestRiwayatPasien(params, data) {
    if (!data.pasiens || data.pasiens.length === 0) {
        return requestAntrean(params);
    }

    const idPasien = data.pasiens[Math.floor(Math.random() * data.pasiens.length)];
    const scenario = Math.random();
    let url = "";

    if (scenario < 0.65) {
        url = `${BASE_URL}/pasien/${idPasien}/riwayat-kunjungan`;
    } else if (scenario < 0.85) {
        url = `${BASE_URL}/pasien/${idPasien}/riwayat-kunjungan?tanggal=${today}&page=1&limit=5`;
    } else {
        url = `${BASE_URL}/pasien/${idPasien}/riwayat-kunjungan?page=2&limit=5`;
    }

    const res = http.get(url, params);
    riwayatPasienDuration.add(res.timings.duration);
    return res;
}

function requestDiagnosaKunjungan(params, data) {
    if (!data.kunjungans || data.kunjungans.length === 0) {
        return requestAntrean(params);
    }
    const idKunjungan = data.kunjungans[Math.floor(Math.random() * data.kunjungans.length)];
    const statusLanjut = statusLanjutList[Math.floor(Math.random() * statusLanjutList.length)];
    const url = `${BASE_URL}/diagnosa/${statusLanjut}/${idKunjungan}`;
    const res = http.get(url, params);
    diagnosaKunjunganDuration.add(res.timings.duration);
    return res;
}

function requestDiagnosaPasien(params, data) {
    if (!data.pasiens || data.pasiens.length === 0) {
        return requestAntrean(params);
    }
    const idPasien = data.pasiens[Math.floor(Math.random() * data.pasiens.length)];
    const statusLanjut = statusLanjutList[Math.floor(Math.random() * statusLanjutList.length)];
    const url = `${BASE_URL}/diagnosa/${statusLanjut}/pasien/${idPasien}`;
    const res = http.get(url, params);
    diagnosaPasienDuration.add(res.timings.duration);
    return res;
}

export default function (data) {
    mixedWorkload(data);
}

export function antreanDokterWorkload(data) {
    const params = {
        headers: {
            "Accept": "application/json",
            "Authorization": `Bearer ${data.token}`,
        },
    };
    const res = requestAntrean(params);
    antreanDokterDuration.add(res.timings.duration);
    const ok = check(res, {
        "status 200": (r) => r.status === 200,
        "success true": (r) => {
            try { return r.json("success") === true; } catch (_) { return false; }
        },
    });
    antreanDokterFail.add(!ok);
    if (!ok && res.status !== 200) {
        console.warn(`[FAIL DOKTER JWT] HTTP ${res.status} on ${res.url}: ${res.body}`);
    }
    sleep(0.05);
}

export function antreanApiKeyWorkload(data) {
    const params = {
        headers: {
            "Accept": "application/json",
            "X-API-Key": data.apiKey || SERVICE_API_KEY,
        },
    };
    const res = requestAntrean(params);
    antreanApiKeyDuration.add(res.timings.duration);
    const ok = check(res, {
        "status 200": (r) => r.status === 200,
        "success true": (r) => {
            try { return r.json("success") === true; } catch (_) { return false; }
        },
    });
    antreanApiKeyFail.add(!ok);
    if (!ok && res.status !== 200) {
        console.warn(`[FAIL API KEY] HTTP ${res.status} on ${res.url}: ${res.body}`);
    }
    sleep(0.05);
}

export function mixedWorkload(data) {
    const authHeaders = (AUTH_MODE === "api_key") ? {
        "Accept": "application/json",
        "X-API-Key": data.apiKey || SERVICE_API_KEY,
    } : {
        "Accept": "application/json",
        "Authorization": `Bearer ${data.token}`,
    };

    const params = {
        headers: authHeaders,
    };

    let res;
    if (ENDPOINT === "obat") {
        res = requestObat(params);
    } else if (ENDPOINT === "icd10") {
        res = requestICD10(params);
    } else if (ENDPOINT === "icd9") {
        res = requestICD9(params);
    } else if (ENDPOINT === "icd") {
        res = Math.random() < 0.6 ? requestICD10(params) : requestICD9(params);
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
    } else if (ENDPOINT === "permintaan_lab_kunjungan") {
        res = requestPermintaanLabKunjungan(params, data);
    } else if (ENDPOINT === "permintaan_lab_pasien") {
        res = requestPermintaanLabPasien(params, data);
    } else if (ENDPOINT === "permintaan_lab_detail") {
        res = requestPermintaanLabDetail(params, data);
    } else if (ENDPOINT === "permintaan_lab") {
        const randP = Math.random();
        if (randP < 0.45) {
            res = requestPermintaanLabKunjungan(params, data);
        } else if (randP < 0.85) {
            res = requestPermintaanLabPasien(params, data);
        } else {
            res = requestPermintaanLabDetail(params, data);
        }
    } else if (ENDPOINT === "laboratorium") {
        const randLab = Math.random();
        if (randLab < 0.25) {
            res = requestLaboratoriumKunjungan(params, data);
        } else if (randLab < 0.50) {
            res = requestLaboratoriumPasien(params, data);
        } else if (randLab < 0.70) {
            res = requestPermintaanLabKunjungan(params, data);
        } else if (randLab < 0.90) {
            res = requestPermintaanLabPasien(params, data);
        } else {
            res = requestPermintaanLabDetail(params, data);
        }
    } else if (ENDPOINT === "rawat_inap") {
        res = requestRawatInap(params, data);
    } else if (ENDPOINT === "riwayat_pasien") {
        res = requestRiwayatPasien(params, data);
    } else if (ENDPOINT === "rawatinap_riwayat") {
        res = Math.random() < 0.5 ? requestRawatInap(params, data) : requestRiwayatPasien(params, data);
    } else if (ENDPOINT === "diagnosa_kunjungan") {
        res = requestDiagnosaKunjungan(params, data);
    } else if (ENDPOINT === "diagnosa_pasien") {
        res = requestDiagnosaPasien(params, data);
    } else if (ENDPOINT === "diagnosa") {
        res = Math.random() < 0.5 ? requestDiagnosaKunjungan(params, data) : requestDiagnosaPasien(params, data);
    } else {
        // Mode 'all': bagi beban ke seluruh modul backend
        const rand = Math.random();
        if (rand < 0.08) {
            res = requestAntrean(params);
        } else if (rand < 0.16) {
            res = requestObat(params);
        } else if (rand < 0.24) {
            res = requestICD10(params);
        } else if (rand < 0.30) {
            res = requestICD9(params);
        } else if (rand < 0.38) {
            res = requestPemeriksaanKunjungan(params, data);
        } else if (rand < 0.46) {
            res = requestPemeriksaanPasien(params, data);
        } else if (rand < 0.54) {
            res = requestResepKunjungan(params, data);
        } else if (rand < 0.62) {
            res = requestResepPasien(params, data);
        } else if (rand < 0.72) {
            res = requestLaboratoriumKunjungan(params, data);
        } else if (rand < 0.82) {
            res = requestPermintaanLabKunjungan(params, data);
        } else if (rand < 0.91) {
            res = requestRawatInap(params, data);
        } else {
            res = requestRiwayatPasien(params, data);
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
// # Test khusus master ICD-10 (In-Memory Cache):
// k6 run -e ENDPOINT=icd10 k6-script.js
//
// # Test khusus master ICD-9 (In-Memory Cache):
// k6 run -e ENDPOINT=icd9 k6-script.js
//
// # Test gabungan ICD-10 & ICD-9 (60:40):
// k6 run -e ENDPOINT=icd k6-script.js
//
// # Quick test ICD (10 detik):
// k6 run -e ENDPOINT=icd -e PROFILE=quick k6-script.js
//
// # Test khusus modul Daftar Pasien Rawat Inap:
// k6 run -e ENDPOINT=rawat_inap k6-script.js
//
// # Test khusus modul Riwayat Kunjungan Pasien (Timeline Index):
// k6 run -e ENDPOINT=riwayat_pasien k6-script.js
//
// # Test gabungan Rawat Inap & Riwayat Pasien (50:50):
// k6 run -e ENDPOINT=rawatinap_riwayat k6-script.js
//
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
