(() => {
  const container = document.getElementById("tesseract-container");
  if (!container) return;

  const scene = new THREE.Scene();
  const camera = new THREE.PerspectiveCamera(
    70,
    container.clientWidth / container.clientHeight,
    0.1,
    1000
  );
  const renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true });
  renderer.setSize(container.clientWidth, container.clientHeight);
  container.innerHTML = "";
  container.appendChild(renderer.domElement);

  camera.position.set(0, 0, 5);
  camera.lookAt(0, 0, 0);

  // === VERTICES 4D ===
  const vertices4D = [];
  for (let x of [-1, 1])
    for (let y of [-1, 1])
      for (let z of [-1, 1])
        for (let w of [-1, 1]) vertices4D.push([x, y, z, w]);

  // === ARÊTES ===
  const edges = [];
  for (let i = 0; i < vertices4D.length; i++) {
    for (let j = i + 1; j < vertices4D.length; j++) {
      const d1 =
        Math.abs(vertices4D[i][0] - vertices4D[j][0]) +
        Math.abs(vertices4D[i][1] - vertices4D[j][1]) +
        Math.abs(vertices4D[i][2] - vertices4D[j][2]) +
        Math.abs(vertices4D[i][3] - vertices4D[j][3]);
      if (d1 === 2) edges.push([i, j]);
    }
  }

  const material = new THREE.LineBasicMaterial({ color: 0x2d79ff });
  const geometry = new THREE.BufferGeometry();
  const mesh = new THREE.LineSegments(geometry, material);
  scene.add(mesh);

  function project4Dto3D([x, y, z, w], ax, ay, az) {
    const cx = Math.cos(ax), sx = Math.sin(ax);
    const cy = Math.cos(ay), sy = Math.sin(ay);
    const cz = Math.cos(az), sz = Math.sin(az);

    let x1 = x * cx - w * sx;
    let w1 = x * sx + w * cx;

    let y1 = y * cy - w1 * sy;
    let w2 = y * sy + w1 * cy;

    let z1 = z * cz - w2 * sz;
    let w3 = z * sz + w2 * cz;

    const d = 3;
    const f = d / (d - w3);
    return new THREE.Vector3(x1 * f, y1 * f, z1 * f);
  }

  // === ANIMATION ===
  let t = 0;
  function animate() {
    requestAnimationFrame(animate);
    t += 0.0045; // ✅ un peu plus rapide

    const ax = t * 0.45, ay = t * 0.55, az = t * 0.4;

    const pulse = (Math.sin(t * 1.3) + 1) * 0.5;
    const lum = 0.55 + pulse * 0.12;
    material.color.setHSL(210 / 360, 1.0, lum);

    camera.position.z = 5 + Math.sin(t * 0.45) * 0.15;
    camera.lookAt(0, 0, 0);

    const points = vertices4D.map(v => project4Dto3D(v, ax, ay, az));

    const positions = new Float32Array(edges.length * 2 * 3);
    let k = 0;
    for (const [i, j] of edges) {
      const a = points[i], b = points[j];
      positions[k++] = a.x; positions[k++] = a.y; positions[k++] = a.z;
      positions[k++] = b.x; positions[k++] = b.y; positions[k++] = b.z;
    }
    geometry.setAttribute("position", new THREE.BufferAttribute(positions, 3));
    geometry.computeBoundingSphere();

    renderer.render(scene, camera);
  }

  function onResize() {
    const w = container.clientWidth || 250;
    const h = container.clientHeight || 250;
    camera.aspect = w / h;
    camera.updateProjectionMatrix();
    renderer.setSize(w, h);
  }
  window.addEventListener("resize", onResize);
  onResize();

  animate();
})();
