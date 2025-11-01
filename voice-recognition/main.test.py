from faster_whisper import WhisperModel
from faster_whisper import BatchedInferencePipeline

model = WhisperModel(
    "large-v3",
    device="cuda",
    compute_type="float16",
    num_workers=4,  # thread CPU per la decodifica
    download_root="./models",  # cache del modello localmente
)
segments, info = model.transcribe("output.wav", language="it")

print(f"Lingua rilevata: {info.language}")
print(f"Durata: {info.duration:.2f}s")

for s in segments:
    print(f"[{s.start:.2f}s -> {s.end:.2f}s] {s.text}")