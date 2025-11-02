import sys
from pathlib import Path

from faster_whisper import BatchedInferencePipeline, WhisperModel
from tqdm import tqdm

if len(sys.argv) < 2:
    print("You need to specify the audio's path.")
    sys.exit()

if len(sys.argv) < 3:
    print("You need to specify the ouput path.")
    sys.exit()

model = WhisperModel(
    sys.argv[3] if len(sys.argv) > 3 else "turbo",
    device="cuda",
    compute_type="float16",
    num_workers=4,
    download_root="./models",  # cache del modello localmente
)

# Crea una pipeline di inferenza in batch
pipe = BatchedInferencePipeline(model)

audioPath = sys.argv[1]
fileName = path = Path(audioPath).stem

outputPath = sys.argv[2]

# Transcrivi l'audio con la pipeline
segments, info = pipe.transcribe(
    audioPath,
    beam_size=5,  # precisione del decoding
    best_of=2,
    vad_filter=True,  # filtra i silenzi
    chunk_length=30,  # dimensione dei chunk in secondi
    batch_size=16,  # numero di chunk elaborati per passaggio
    language=sys.argv[4] if len(sys.argv) > 4 else None,  # lingua del testo trascritto
)

segments = list(segments)

print(f"Language: {info.language}")
print(f"Duration: {info.duration:.2f}s")
print(f"Segments: {len(list(segments))}")

completeOutput = outputPath + "\\" + fileName + ".txt"
print(f"Transcribing {audioPath} to {completeOutput}")

with open(completeOutput, "w+", encoding="utf-8") as f:
    for segment in tqdm(segments, desc="Processing"):
        f.write(segment.text.strip() + " ")
