# -*- coding: utf-8 -*-
"""
美式鄉村風 客廳 + 餐廳 — Blender 程序化場景建模腳本
=====================================================
使用方式:
  1. 開啟 Blender (3.6 / 4.x)
  2. Scripting 工作區 → New → 貼上本檔 → Run Script (Alt+P)
     或命令列: blender --python american_country_room.py
  3. 注意: 會清空目前檔案中的所有物件 / 集合

座標: X = 房間長邊 (0 → 10m, 左 = 客廳, 右 = 餐廳/廚房)
      Y = 深度 (0 = 前方開口 / 攝影機側, 6m = 後牆)
      Z = 高度 (天花板 2.9m)

場景內容 (對照手繪圖):
  - 淺藍牆面 + 白色腰牆飾板 (wainscoting) + 踢腳板 + 冠頂線板
  - 客廳跌級天花 (tray ceiling) + 嵌燈; 餐廳木樑天花; 廚房方形天窗
  - 白色文化石電視牆 + 拱形壁龕 + 裝飾盤 + 爐台
  - 白色書櫃 (下櫃門片 + 層板 + 隨機書本/花瓶/盤子)
  - 藍色捲手沙發、米色/黃色翼背扶手椅、車腳咖啡桌、圓邊几、檯燈
  - 長餐桌 + 8 張溫莎椅 + 鐵件圓環燭台吊燈
  - 客廳燈罩式水晶吊燈、吧台中島 + 吊櫃 + 酒杯架 + 吊燈
  - 格子窗、拱形窗、拱形雙開法式落地門、鑲板門
  - 盆栽植物、地毯、掛畫
  - 兩台攝影機、Cycles 設定、可選 Freestyle 線稿 (模擬手繪感)
"""

import bpy
import bmesh
import random
import time
from math import pi, sin, cos, sqrt, radians, atan2
from mathutils import Vector

random.seed(27)

# ---------------------------------------------------------------- 參數
W, D, H, T = 10.0, 6.0, 2.9, 0.15     # 房間 寬 / 深 / 高 / 牆厚
SKETCH_LINES = True                   # Freestyle 線稿, 模擬手繪透視圖
SAMPLES = 128


# ================================================================ 基礎工具
_COLL = None


def use_coll(name):
    """之後建立的物件都放進這個 collection"""
    global _COLL
    c = bpy.data.collections.get(name) or bpy.data.collections.new(name)
    root = bpy.context.scene.collection
    if root.children.get(c.name) is None:
        root.children.link(c)
    _COLL = c


def _link(o, parent=None):
    (_COLL or bpy.context.scene.collection).objects.link(o)
    if parent is not None:
        o.parent = parent
    return o


def _obj_from_bm(name, bm, mat, parent, loc, rot, smooth=False, bevel=0.0):
    if smooth == 'SIDES':
        bm.normal_update()
        for f in bm.faces:
            f.smooth = abs(f.normal.z) < 0.9
    elif smooth:
        for f in bm.faces:
            f.smooth = True
    me = bpy.data.meshes.new(name)
    bm.to_mesh(me)
    bm.free()
    if mat is not None:
        me.materials.append(mat)
    o = bpy.data.objects.new(name, me)
    o.location = loc
    o.rotation_euler = rot
    if bevel > 0:
        m = o.modifiers.new("Bevel", 'BEVEL')
        m.width = bevel
        m.segments = 3
        m.limit_method = 'ANGLE'
    return _link(o, parent)


def box(name, c, s, mat=None, parent=None, rot=(0, 0, 0), bevel=0.0):
    bm = bmesh.new()
    bmesh.ops.create_cube(bm, size=1.0)
    bmesh.ops.scale(bm, vec=Vector(s), verts=bm.verts)
    return _obj_from_bm(name, bm, mat, parent, c, rot, bevel=bevel)


def cyl(name, c, r, h, mat=None, parent=None, rot=(0, 0, 0), r2=None, seg=24, caps=True, bevel=0.0):
    """圓柱 / 圓錐 (r = 底部半徑, r2 = 頂部半徑)"""
    bm = bmesh.new()
    bmesh.ops.create_cone(bm, cap_ends=caps, cap_tris=False, segments=seg,
                          radius1=r, radius2=r if r2 is None else r2, depth=h)
    return _obj_from_bm(name, bm, mat, parent, c, rot, smooth='SIDES', bevel=bevel)


def sphere(name, c, r, mat=None, parent=None, scale=(1, 1, 1), rot=(0, 0, 0), seg=20):
    bm = bmesh.new()
    bmesh.ops.create_uvsphere(bm, u_segments=seg, v_segments=max(6, seg // 2), radius=r)
    bmesh.ops.scale(bm, vec=Vector(scale), verts=bm.verts)
    return _obj_from_bm(name, bm, mat, parent, c, rot, smooth=True)


def rod(name, p1, p2, r, mat=None, parent=None, r2=None, seg=10):
    """兩點之間的圓柱 (r 在 p1 端, r2 在 p2 端)"""
    p1, p2 = Vector(p1), Vector(p2)
    d = p2 - p1
    q = Vector((0, 0, 1)).rotation_difference(d)
    return cyl(name, (p1 + p2) / 2, r, d.length, mat, parent, rot=q.to_euler(), r2=r2, seg=seg)


def torus(name, c, R, r, mat=None, parent=None, rot=(0, 0, 0), arc=2 * pi, seg=40, rseg=8):
    """圓環; arc < 2π 時為一段弧 (從 +X 往 +Y 方向)"""
    bm = bmesh.new()
    closed = arc >= 2 * pi - 1e-6
    n_rings = seg if closed else seg + 1
    rings = []
    for i in range(n_rings):
        a = arc * i / seg
        rings.append([bm.verts.new(((R + r * cos(b)) * cos(a), (R + r * cos(b)) * sin(a), r * sin(b)))
                      for b in (2 * pi * j / rseg for j in range(rseg))])
    for i in range(seg):
        r0, r1 = rings[i], rings[(i + 1) % n_rings]
        for j in range(rseg):
            k = (j + 1) % rseg
            bm.faces.new((r0[j], r1[j], r1[k], r0[k]))
    bmesh.ops.recalc_face_normals(bm, faces=bm.faces[:])
    return _obj_from_bm(name, bm, mat, parent, c, rot, smooth=True)


def mesh(name, verts, faces, mat=None, parent=None):
    bm = bmesh.new()
    vs = [bm.verts.new(v) for v in verts]
    for f in faces:
        bm.faces.new([vs[i] for i in f])
    return _obj_from_bm(name, bm, mat, parent, (0, 0, 0), (0, 0, 0))


def empty(name, loc=(0, 0, 0), rot_z=0.0, parent=None):
    o = bpy.data.objects.new(name, None)
    o.empty_display_size = 0.25
    o.location = loc
    o.rotation_euler = (0, 0, rot_z)
    return _link(o, parent)


def facing(dx, dy):
    """家具本地座標正面朝 -Y; 回傳讓正面朝向 (dx, dy) 的 Z 旋轉角"""
    return atan2(dx, -dy)


def light(name, loc, kind='POINT', power=10.0, color=(1.0, 0.85, 0.65), size=0.03, parent=None, rot=(0, 0, 0)):
    ld = bpy.data.lights.new(name, kind)
    ld.energy = power
    ld.color = color
    if kind in ('POINT', 'SPOT'):
        ld.shadow_soft_size = size
    if kind == 'SPOT':
        ld.spot_size = radians(110)
        ld.spot_blend = 0.6
    if kind == 'AREA':
        ld.shape = 'RECTANGLE'
        ld.size, ld.size_y = size
    if kind == 'SUN':
        ld.angle = radians(3)
    o = bpy.data.objects.new(name, ld)
    o.location = loc
    o.rotation_euler = rot
    return _link(o, parent)


def boolean_cut(target, cutters):
    """用 cutters 對 target 做差集並套用 (不依賴 bpy.ops)"""
    for c in cutters:
        m = target.modifiers.new("Cut", 'BOOLEAN')
        m.operation = 'DIFFERENCE'
        m.object = c
        try:
            m.solver = 'EXACT'
        except Exception:
            pass
    bpy.context.view_layer.update()
    dg = bpy.context.evaluated_depsgraph_get()
    new_me = bpy.data.meshes.new_from_object(target.evaluated_get(dg))
    old = target.data
    target.modifiers.clear()
    target.data = new_me
    new_me.name = target.name
    bpy.data.meshes.remove(old)
    for c in cutters:
        bpy.data.objects.remove(c, do_unlink=True)


def molding_rect(name, parent, x0, x1, z0, z1, y, mat, t=0.025, d=0.02):
    """XZ 平面上的矩形線板框 (門片 / 腰牆飾板)"""
    xm, zm = (x0 + x1) / 2, (z0 + z1) / 2
    box(name + "_B", (xm, y, z0), (x1 - x0 + t, d, t), mat, parent)
    box(name + "_T", (xm, y, z1), (x1 - x0 + t, d, t), mat, parent)
    box(name + "_L", (x0, y, zm), (t, d, z1 - z0), mat, parent)
    box(name + "_R", (x1, y, zm), (t, d, z1 - z0), mat, parent)


def rect_ring(name, x0, x1, y0, y1, wd, zc, h, mat, parent=None):
    """XY 平面上的矩形框 (天花板跌級 / 樑框)"""
    box(name + "_S", ((x0 + x1) / 2, y0 + wd / 2, zc), (x1 - x0, wd, h), mat, parent)
    box(name + "_N", ((x0 + x1) / 2, y1 - wd / 2, zc), (x1 - x0, wd, h), mat, parent)
    box(name + "_W", (x0 + wd / 2, (y0 + y1) / 2, zc), (wd, y1 - y0 - 2 * wd, h), mat, parent)
    box(name + "_E", (x1 - wd / 2, (y0 + y1) / 2, zc), (wd, y1 - y0 - 2 * wd, h), mat, parent)


def arc_boxes(name, parent, cx, zs, R, thick, depth, y, mat, n=14):
    """以小方塊排出 XZ 平面上的半圓拱"""
    for i in range(n):
        a0, a1 = pi * i / n, pi * (i + 1) / n
        am = (a0 + a1) / 2
        seg = 2 * R * sin((a1 - a0) / 2) * 1.08
        box(f"{name}{i}", (cx + R * cos(am), y, zs + R * sin(am)), (seg, depth, thick),
            mat, parent, rot=(0, pi / 2 - am, 0))


# ================================================================ 材質
def _principled(name):
    m = bpy.data.materials.new(name)
    try:
        m.use_nodes = True
    except Exception:
        pass
    nt = m.node_tree
    b = next((n for n in nt.nodes if n.type == 'BSDF_PRINCIPLED'), None)
    if b is None:
        nt.nodes.clear()
        out = nt.nodes.new("ShaderNodeOutputMaterial")
        b = nt.nodes.new("ShaderNodeBsdfPrincipled")
        nt.links.new(b.outputs[0], out.inputs[0])
    return m, b


def _set(node, names, val):
    for n in ((names,) if isinstance(names, str) else names):
        if n in node.inputs:
            node.inputs[n].default_value = val
            return


def mat(name, color, rough=0.5, metal=0.0, emit=None, strength=0.0, alpha=1.0, sheen=0.0):
    m, b = _principled(name)
    _set(b, "Base Color", (*color, 1))
    _set(b, "Roughness", rough)
    _set(b, "Metallic", metal)
    if emit:
        _set(b, ("Emission Color", "Emission"), (*emit, 1))
        _set(b, "Emission Strength", strength)
    if sheen:
        _set(b, ("Sheen Weight", "Sheen"), sheen)
    if alpha < 1:
        _set(b, "Alpha", alpha)
        for attr, v in (("blend_method", 'BLEND'), ("surface_render_method", 'BLENDED')):
            try:
                setattr(m, attr, v)
            except Exception:
                pass
    m.diffuse_color = (*color, alpha)
    return m


def mat_wood(name, c1, c2, scale=5.0, rough=0.45):
    """程序化木紋 (Wave Texture)"""
    m, b = _principled(name)
    N, L = m.node_tree.nodes, m.node_tree.links
    tc = N.new("ShaderNodeTexCoord")
    wv = N.new("ShaderNodeTexWave")
    wv.wave_type = 'BANDS'
    try:
        wv.bands_direction = 'X'
    except Exception:
        pass
    _set(wv, "Scale", scale)
    _set(wv, "Distortion", 6.0)
    _set(wv, "Detail", 4.0)
    ramp = N.new("ShaderNodeValToRGB")
    ramp.color_ramp.elements[0].color = (*c1, 1)
    ramp.color_ramp.elements[1].color = (*c2, 1)
    L.new(tc.outputs["Object"], wv.inputs["Vector"])
    L.new(wv.outputs["Fac"], ramp.inputs["Fac"])
    L.new(ramp.outputs["Color"], b.inputs["Base Color"])
    _set(b, "Roughness", rough)
    m.diffuse_color = (*c2, 1)
    return m


def mat_boards(name, c1, c2, mortar, brick_w, row_h, gap, rough=0.5, vertical=False, offset=0.5, bump=0.3):
    """用 Brick Texture 做 地板木條 / 文化石 / 橫向企口板"""
    m, b = _principled(name)
    N, L = m.node_tree.nodes, m.node_tree.links
    tc = N.new("ShaderNodeTexCoord")
    br = N.new("ShaderNodeTexBrick")
    br.offset = offset
    _set(br, "Color1", (*c1, 1))
    _set(br, "Color2", (*c2, 1))
    _set(br, "Mortar", (*mortar, 1))
    _set(br, "Scale", 1.0)
    _set(br, "Mortar Size", gap)
    _set(br, "Mortar Smooth", 0.1)
    _set(br, "Brick Width", brick_w)
    _set(br, "Row Height", row_h)
    src = tc.outputs["Object"]
    if vertical:  # 貼在 XZ 牆面上
        mp = N.new("ShaderNodeMapping")
        _set(mp, "Rotation", (pi / 2, 0, 0))
        L.new(src, mp.inputs["Vector"])
        src = mp.outputs["Vector"]
    L.new(src, br.inputs["Vector"])
    L.new(br.outputs["Color"], b.inputs["Base Color"])
    bp = N.new("ShaderNodeBump")
    bp.invert = True
    _set(bp, "Strength", bump)
    _set(bp, "Distance", 0.02)
    L.new(br.outputs["Fac"], bp.inputs["Height"])
    L.new(bp.outputs["Normal"], b.inputs["Normal"])
    _set(b, "Roughness", rough)
    m.diffuse_color = (*c1, 1)
    return m


def build_materials():
    M = {
        'wall': mat("Wall_Blue", (0.42, 0.62, 0.85), 0.85),
        'trim': mat("Trim_White", (0.9, 0.9, 0.86), 0.35),
        'ceiling': mat("Ceiling_White", (0.93, 0.93, 0.9), 0.9),
        'floor': mat_boards("Floor_Planks", (0.55, 0.38, 0.22), (0.68, 0.5, 0.32), (0.25, 0.17, 0.1),
                            1.4, 0.13, 0.002, rough=0.35, offset=0.37),
        'brick': mat_boards("Brick_WhiteWash", (0.86, 0.83, 0.78), (0.74, 0.68, 0.6), (0.95, 0.94, 0.92),
                            0.23, 0.075, 0.008, rough=0.9, vertical=True),
        'shiplap': mat_boards("Shiplap_White", (0.9, 0.9, 0.86), (0.86, 0.86, 0.82), (0.5, 0.5, 0.48),
                              8.0, 0.12, 0.004, rough=0.5, vertical=True),
        'ceil_wood': mat_boards("Ceiling_Planks", (0.72, 0.55, 0.35), (0.8, 0.63, 0.42), (0.35, 0.25, 0.15),
                                2.0, 0.15, 0.004, rough=0.6),
        'beam': mat_wood("Wood_Beam", (0.5, 0.33, 0.18), (0.62, 0.43, 0.25), scale=3),
        'wood_dark': mat_wood("Wood_Dark", (0.18, 0.09, 0.04), (0.3, 0.16, 0.07)),
        'wood_mid': mat_wood("Wood_Mid", (0.4, 0.22, 0.1), (0.55, 0.33, 0.16)),
        'wood_light': mat_wood("Wood_Light", (0.6, 0.45, 0.28), (0.72, 0.56, 0.36)),
        'sofa': mat("Fabric_Blue", (0.05, 0.3, 0.6), 0.85, sheen=0.5),
        'cream': mat("Fabric_Cream", (0.82, 0.75, 0.6), 0.9, sheen=0.3),
        'yellow': mat("Fabric_Yellow", (0.85, 0.68, 0.15), 0.9, sheen=0.3),
        'pillow': mat("Fabric_Pillow", (0.9, 0.85, 0.7), 0.9),
        'iron': mat("Iron_Black", (0.03, 0.03, 0.03), 0.45, 0.7),
        'brass': mat("Brass", (0.8, 0.58, 0.28), 0.3, 1.0),
        'glass': mat("Glass", (0.85, 0.95, 1.0), 0.02, alpha=0.12),
        'shade': mat("Lamp_Shade", (0.95, 0.9, 0.78), 0.8, emit=(1.0, 0.85, 0.6), strength=1.5),
        'candle': mat("Candle_Wax", (0.95, 0.93, 0.88), 0.5),
        'flame': mat("Flame", (1.0, 0.7, 0.3), 0.5, emit=(1.0, 0.7, 0.3), strength=20),
        'downlight': mat("Downlight_Lens", (1, 1, 1), 0.3, emit=(1.0, 0.95, 0.85), strength=15),
        'skylight': mat("Skylight_Panel", (1, 1, 1), 0.3, emit=(0.9, 0.95, 1.0), strength=4),
        'tv': mat("TV_Black", (0.01, 0.01, 0.012), 0.1),
        'stone': mat("Stone_Hearth", (0.7, 0.66, 0.6), 0.8),
        'terracotta': mat("Terracotta", (0.62, 0.3, 0.16), 0.7),
        'ceramic_white': mat("Ceramic_White", (0.92, 0.92, 0.9), 0.15),
        'ceramic_blue': mat("Ceramic_Blue", (0.15, 0.3, 0.6), 0.15),
        'leaf': mat("Leaf_Dark", (0.13, 0.38, 0.1), 0.5),
        'leaf2': mat("Leaf_Light", (0.25, 0.5, 0.15), 0.5),
        'stem': mat("Stem", (0.2, 0.3, 0.1), 0.6),
        'soil': mat("Soil", (0.12, 0.08, 0.05), 1.0),
        'rug': mat("Rug", (0.78, 0.7, 0.55), 1.0),
        'rug_border': mat("Rug_Border", (0.55, 0.42, 0.3), 1.0),
        'ground': mat("Ground", (0.35, 0.45, 0.25), 1.0),
        'flower': mat("Flower", (0.9, 0.4, 0.5), 0.6),
        'fruit': mat("Fruit", (0.85, 0.4, 0.1), 0.4),
        'art1': mat("Art_1", (0.35, 0.55, 0.45), 0.6),
        'art2': mat("Art_2", (0.8, 0.65, 0.4), 0.6),
    }
    book_cols = [(0.55, 0.12, 0.1), (0.1, 0.25, 0.45), (0.15, 0.35, 0.2),
                 (0.75, 0.65, 0.45), (0.3, 0.2, 0.15), (0.85, 0.82, 0.75)]
    M['books'] = [mat(f"Book_{i}", c, 0.6) for i, c in enumerate(book_cols)]
    return M


# ================================================================ 牆面 / 門窗
def opening_unit(name, e, op, M):
    """格子窗 / 拱窗 / 法式門 (在牆本地座標建立, 室內為 -Y)"""
    x, w, z0, zs = op['x'], op['w'], op['z0'], op['z1']
    arch, leaves = op.get('arch', False), op.get('leaves', 1)
    cols, rows, kick = op.get('cols', 2), op.get('rows', 3), op.get('kick', 0.0)
    fr, gl = M['trim'], M['glass']
    fw, mw = 0.055, 0.022
    yc, dp = T / 2, T * 0.55
    xl, xr = x - w / 2, x + w / 2
    R = w / 2
    zg = z0 + kick
    zm = (z0 + zs) / 2

    # 玻璃
    box(name + "_Glass", (x, yc, (zg + zs) / 2), (w, 0.006, zs - zg), gl, e)
    if arch:
        n = 32
        verts = [(x, yc, zs)] + [(x + R * cos(pi * i / n), yc, zs + R * sin(pi * i / n)) for i in range(n + 1)]
        mesh(name + "_GlassArch", verts, [(0, i + 1, i + 2) for i in range(n)], gl, e)

    # 框
    box(name + "_StileL", (xl + fw / 2, yc, zm), (fw, dp, zs - z0), fr, e)
    box(name + "_StileR", (xr - fw / 2, yc, zm), (fw, dp, zs - z0), fr, e)
    box(name + "_RailB", (x, yc, z0 + fw / 2), (w, dp, fw), fr, e)
    box(name + "_RailT", (x, yc, zs - (0 if arch else fw / 2)), (w, dp, fw), fr, e)
    top = zs - (fw / 2 if arch else fw)
    lw = w / leaves
    for li in range(leaves):
        a, b = xl + li * lw, xl + (li + 1) * lw
        if li > 0:
            box(f"{name}_Meet{li}", (a, yc, zm), (fw * 1.4, dp, zs - z0), fr, e)
        if kick > 0:   # 法式門下半部實心門板
            box(f"{name}_Kick{li}", ((a + b) / 2, yc, z0 + kick / 2), (lw - fw, dp * 0.8, kick), fr, e)
            box(f"{name}_KickRail{li}", ((a + b) / 2, yc, zg), (lw, dp, fw), fr, e)
            molding_rect(f"{name}_KickPanel{li}", e, a + 0.1, b - 0.1, z0 + 0.1, zg - 0.08,
                         yc - dp * 0.4 - 0.005, fr, t=0.02, d=0.01)
        for c in range(1, cols):   # 格子窗櫺
            mx = a + (b - a) * c / cols
            box(f"{name}_MV{li}_{c}", (mx, yc, (zg + top) / 2), (mw, dp * 0.7, top - zg), fr, e)
        for r_ in range(1, rows):
            mz = zg + (top - zg) * r_ / rows
            box(f"{name}_MH{li}_{r_}", ((a + b) / 2, yc, mz), (lw, dp * 0.7, mw), fr, e)

    if arch:   # 拱頂框 + 放射窗櫺
        arc_boxes(name + "_ArchFr", e, x, zs, R - fw / 2, fw, dp, yc, fr, n=16)
        for k, a in enumerate((pi / 4, pi / 2, 3 * pi / 4)):
            box(f"{name}_Spoke{k}", (x + R / 2 * cos(a), yc, zs + R / 2 * sin(a)),
                (R, dp * 0.7, mw), fr, e, rot=(0, -a, 0))
        arc_boxes(name + "_Hub", e, x, zs, R * 0.3, mw, dp * 0.7, yc, fr, n=8)

    # 室內側門窗套
    cy, cd, cw = -0.015, 0.03, 0.08
    box(name + "_CasL", (xl - cw / 2, cy, zm), (cw, cd, zs - z0), fr, e)
    box(name + "_CasR", (xr + cw / 2, cy, zm), (cw, cd, zs - z0), fr, e)
    if arch:
        arc_boxes(name + "_CasArch", e, x, zs, R + cw / 2, cw, cd, cy, fr, n=20)
        box(name + "_Keystone", (x, cy - 0.01, zs + R + cw * 0.6), (0.1, cd + 0.02, 0.16), fr, e)
    else:
        box(name + "_CasHead", (x, cy - 0.005, zs + 0.06), (w + 2 * cw + 0.06, cd + 0.01, 0.12), fr, e, bevel=0.005)
    if z0 > 0.01:
        box(name + "_Sill", (x, (T - 0.1) / 2, z0 - 0.02), (w + 2 * cw + 0.04, T + 0.1, 0.04), fr, e, bevel=0.005)


def wainscot(name, e, length, openings, blocks, M):
    """腰牆飾板 + 腰線 + 踢腳板, 避開門窗與櫃體"""
    tr = M['trim']
    blocked = [(op['x'] - op['w'] / 2 - 0.1, op['x'] + op['w'] / 2 + 0.1) for op in openings if op['z0'] < 0.85]
    blocked = sorted(blocked + list(blocks))
    spans, cur = [], 0.0
    for a, b in blocked:
        if a > cur:
            spans.append((cur, a))
        cur = max(cur, b)
    if cur < length:
        spans.append((cur, length))
    for k, (a, b) in enumerate(spans):
        L = b - a
        if L < 0.12:
            continue
        cx = (a + b) / 2
        box(f"{name}_WBoard{k}", (cx, -0.006, 0.45), (L, 0.012, 0.9), tr, e)
        box(f"{name}_WRail{k}", (cx, -0.025, 0.9), (L, 0.05, 0.05), tr, e, bevel=0.01)
        box(f"{name}_WBase{k}", (cx, -0.015, 0.07), (L, 0.03, 0.14), tr, e)
        n = max(1, int(L / 0.75))
        gap = 0.09
        pw = (L - gap * (n + 1)) / n
        if pw < 0.12:
            continue
        for i in range(n):
            x0 = a + gap + i * (pw + gap)
            molding_rect(f"{name}_WPanel{k}_{i}", e, x0, x0 + pw, 0.24, 0.76, -0.02, tr)


def wall_unit(name, origin, rot_z, length, openings, blocks, M):
    """一面牆: 本地 X 沿牆, 牆面在 y=0, 室內 = -Y"""
    use_coll("Room")
    e = empty(name, origin, rot_z)
    wall = box(name + "_Wall", (length / 2, T / 2, H / 2), (length + 2 * T, T, H), M['wall'], e)
    cutters = []
    for i, op in enumerate(openings):
        z0 = op['z0'] if op['z0'] > 0 else -0.05
        zs = op['z1']
        cutters.append(box(f"{name}_cut{i}", (op['x'], T / 2, (z0 + zs) / 2), (op['w'], T * 4, zs - z0), None, e))
        if op.get('arch'):
            cutters.append(cyl(f"{name}_cutA{i}", (op['x'], T / 2, zs), op['w'] / 2, T * 4, None, e,
                               rot=(pi / 2, 0, 0), seg=48))
    if cutters:
        boolean_cut(wall, cutters)
    use_coll("Windows_Doors")
    for i, op in enumerate(openings):
        opening_unit(f"{name}_Open{i}", e, op, M)
    use_coll("Trim")
    wainscot(name, e, length, openings, blocks, M)
    box(name + "_Crown", (length / 2, -0.04, H - 0.06), (length, 0.08, 0.12), M['trim'], e, bevel=0.01)
    return e


def wall_frame(name, parent, x, z, w, h, art, M):
    box(name + "_Frame", (x, -0.0125, z), (w, 0.025, h), M['wood_dark'], parent, bevel=0.004)
    box(name + "_Art", (x, -0.026, z), (w - 0.06, 0.004, h - 0.06), art, parent)


def build_room(M):
    use_coll("Room")
    box("Floor", (W / 2, D / 2, -0.025), (W, D, 0.05), M['floor'])
    box("Ground_Outside", (W / 2, D / 2, -0.08), (60, 60, 0.04), M['ground'])
    box("Ceiling", (W / 2, D / 2, H + 0.05), (W + 2 * T, D + 2 * T, 0.1), M['ceiling'])

    # 後牆: 餐廳三樘格子窗; 電視牆與書櫃區不做腰牆
    back_ops = [dict(x=x, w=0.8, z0=0.9, z1=2.3, cols=2, rows=3) for x in (7.05, 8.05, 9.05)]
    wall_unit("BackWall", (0, D, 0), 0.0, W, back_ops, [(1.15, 3.85), (4.05, 6.35)], M)

    # 左牆: 拱形窗 + 鑲板門
    left_ops = [dict(x=1.6, w=0.9, z0=0.9, z1=1.85, arch=True, cols=2, rows=3)]
    left = wall_unit("LeftWall", (0, 0, 0), pi / 2, D, left_ops, [(4.2, 5.4)], M)

    # 右牆: 拱形雙開法式落地門 (本地 x = D - 世界 y)
    right_ops = [dict(x=D - 1.6, w=1.3, z0=0.0, z1=2.0, arch=True, leaves=2, cols=2, rows=4, kick=0.45)]
    right = wall_unit("RightWall", (W, D, 0), -pi / 2, D, right_ops, [], M)

    # 左牆鑲板門
    use_coll("Windows_Doors")
    tr = M['trim']
    dx, dw, dh = 4.8, 0.9, 2.1
    box("LeftDoor_Slab", (dx, -0.02, dh / 2), (dw, 0.04, dh), tr, left)
    for k, (z0, z1) in enumerate(((0.15, 0.95), (1.1, 1.95))):
        molding_rect(f"LeftDoor_Panel{k}", left, dx - dw / 2 + 0.12, dx + dw / 2 - 0.12, z0, z1, -0.045, tr)
    for s in (-1, 1):
        box(f"LeftDoor_Cas{s}", (dx + s * (dw / 2 + 0.04), -0.015, dh / 2), (0.08, 0.03, dh), tr, left)
    box("LeftDoor_Head", (dx, -0.02, dh + 0.06), (dw + 0.22, 0.04, 0.12), tr, left, bevel=0.005)
    sphere("LeftDoor_Knob", (dx + dw / 2 - 0.09, -0.07, 1.0), 0.03, M['brass'], left)

    # 掛畫
    use_coll("Decor")
    wall_frame("Art_L1", left, 2.85, 1.55, 0.4, 0.5, M['art1'], M)
    wall_frame("Art_L2", left, 3.45, 1.65, 0.35, 0.35, M['art2'], M)
    wall_frame("Art_R1", right, 0.9, 1.6, 0.5, 0.6, M['art2'], M)
    return left, right


# ================================================================ 天花板
def downlight(name, x, y, z, M, power=25):
    cyl(name + "_Trim", (x, y, z - 0.004), 0.07, 0.008, M['trim'])
    cyl(name + "_Lens", (x, y, z - 0.009), 0.05, 0.004, M['downlight'])
    light(name + "_Spot", (x, y, z - 0.02), 'SPOT', power, (1.0, 0.88, 0.72), size=0.04)


def build_ceiling(M):
    use_coll("Ceiling")
    c, tr = M['ceiling'], M['trim']

    # 客廳 跌級天花
    lx1, sw = 5.8, 0.55
    rect_ring("Tray_Soffit", 0, lx1, 0, D, sw, H - 0.11, 0.22, c)
    rect_ring("Tray_Lip", sw - 0.03, lx1 - sw + 0.03, sw - 0.03, D - sw + 0.03, 0.06, H - 0.235, 0.05, tr)
    rect_ring("Tray_Crown", sw, lx1 - sw, sw, D - sw, 0.05, H - 0.03, 0.06, tr)
    rect_ring("Tray_Inner", sw + 0.35, lx1 - sw - 0.35, sw + 0.35, D - sw - 0.35, 0.04, H - 0.01, 0.02, tr)

    # 餐廳 木樑天花
    bx0, bx1, by0, by1 = 6.1, W - 0.1, 2.9, D - 0.1
    box("Dining_PlankCeiling", ((bx0 + bx1) / 2, (by0 + by1) / 2, H - 0.006), (bx1 - bx0, by1 - by0, 0.012), M['ceil_wood'])
    rect_ring("Dining_BeamFrame", bx0, bx1, by0, by1, 0.16, H - 0.1, 0.2, M['beam'])
    nb = 5
    for i in range(nb):
        x = bx0 + (bx1 - bx0) * (i + 1) / (nb + 1)
        box(f"Dining_Beam{i}", (x, (by0 + by1) / 2, H - 0.08), (0.12, by1 - by0 - 0.3, 0.16), M['beam'], bevel=0.01)
    box("Dining_CrossBeam", ((bx0 + bx1) / 2, (by0 + by1) / 2, H - 0.06), (bx1 - bx0 - 0.3, 0.1, 0.12), M['beam'])

    # 廚房 方形天窗燈箱
    sx0, sx1, sy0, sy1 = 7.0, 8.8, 0.6, 1.9
    rect_ring("Skylight_Frame", sx0, sx1, sy0, sy1, 0.1, H - 0.06, 0.12, tr)
    rect_ring("Skylight_Frame2", sx0 + 0.1, sx1 - 0.1, sy0 + 0.1, sy1 - 0.1, 0.04, H - 0.02, 0.04, tr)
    box("Skylight_Panel", ((sx0 + sx1) / 2, (sy0 + sy1) / 2, H - 0.003), (sx1 - sx0 - 0.2, sy1 - sy0 - 0.2, 0.004), M['skylight'])

    # 嵌燈
    use_coll("Lights")
    light("Skylight_Area", ((sx0 + sx1) / 2, (sy0 + sy1) / 2, H - 0.02), 'AREA', 120, (0.95, 0.97, 1.0),
          size=(sx1 - sx0 - 0.2, sy1 - sy0 - 0.2))
    zs = H - 0.22
    spots = [(1.7, D - 0.48, zs), (3.3, D - 0.48, zs), (4.6, D - 0.48, zs), (5.5, D - 0.48, zs),
             (0.275, 1.0, zs), (0.275, 2.6, zs), (0.275, 4.2, zs),
             (1.8, 0.275, zs), (3.6, 0.275, zs), (5.2, 0.275, zs),
             (6.4, 1.2, H), (6.4, 2.4, H), (9.5, 2.4, H)]
    for i, (x, y, z) in enumerate(spots):
        downlight(f"Downlight{i}", x, y, z, M)


# ================================================================ 電視牆 (文化石 + 拱形壁龕)
def build_tv_wall(M):
    use_coll("TV_Wall")
    tr = M['trim']
    x0, x1 = 1.2, 3.8
    cx, fy = (x0 + x1) / 2, D - 0.4
    body = box("TVWall_Brick", (cx, D - 0.2, H / 2), (x1 - x0, 0.4, H), M['brick'])
    nw, nz0, nzs = 1.5, 0.55, 1.75
    nr = nw / 2
    boolean_cut(body, [box("TV_cut1", (cx, fy, (nz0 + nzs) / 2), (nw, 0.3, nzs - nz0)),
                       cyl("TV_cut2", (cx, fy, nzs), nr, 0.3, rot=(pi / 2, 0, 0), seg=48)])
    box("TV_NicheBack", (cx, fy + 0.15, (nz0 + nzs + nr) / 2), (nw + 0.1, 0.004, nzs + nr - nz0 + 0.1), M['shiplap'])

    arc_boxes("TV_ArchTrim", None, cx, nzs, nr + 0.035, 0.07, 0.03, fy - 0.01, tr, n=20)
    for s in (-1, 1):
        box(f"TV_Jamb{s}", (cx + s * (nr + 0.035), fy - 0.01, (nz0 + nzs) / 2), (0.07, 0.03, nzs - nz0), tr)
        box(f"TV_Pilaster{s}", (cx + s * ((x1 - x0) / 2 - 0.07), fy - 0.03, H / 2), (0.14, 0.06, H), tr)
    box("TV_NicheSill", (cx, fy + 0.06, nz0 + 0.02), (nw + 0.14, 0.2, 0.04), tr, bevel=0.005)
    box("TV_Cornice", (cx, fy - 0.05, H - 0.1), (x1 - x0 + 0.1, 0.1, 0.2), tr, bevel=0.01)

    box("TV_Screen", (cx, fy + 0.12, 1.12), (1.12, 0.035, 0.64), M['tv'], bevel=0.004)
    plates = [(-0.42, 1.62, 0.075), (-0.22, 1.9, 0.1), (0.0, 2.15, 0.12), (0.22, 1.9, 0.1), (0.42, 1.62, 0.075)]
    for i, (dx, z, r) in enumerate(plates):
        cyl(f"TV_Plate{i}", (cx + dx, fy + 0.135, z), r, 0.012, M['ceramic_white'], rot=(pi / 2, 0, 0), seg=32)
        cyl(f"TV_PlateInner{i}", (cx + dx, fy + 0.128, z), r * 0.65, 0.004, M['ceramic_blue'], rot=(pi / 2, 0, 0), seg=32)

    box("TV_Hearth", (cx, fy - 0.17, 0.2), (x1 - x0 + 0.3, 0.34, 0.4), M['stone'], bevel=0.015)
    cyl("Hearth_Jug", (x0 + 0.2, fy - 0.17, 0.55), 0.08, 0.3, M['terracotta'], r2=0.05)
    sphere("Hearth_Jar", (x0 + 0.42, fy - 0.15, 0.48), 0.08, M['ceramic_blue'], scale=(1, 1, 1.1))
    potted_plant("Hearth_Plant", (x1 - 0.2, fy - 0.17, 0.4), M, pot_r=0.1, pot_h=0.16, n=9,
                 leaf_len=0.18, leaf_w=0.05, height=0.25, pot='ceramic_white')


# ================================================================ 書櫃
def fill_shelf(prefix, x0, x1, z, yc, M):
    x, k = x0 + 0.03, 0
    while x < x1 - 0.06:
        r = random.random()
        if r < 0.55:      # 一排書
            for _ in range(random.randint(4, 10)):
                bw, bh, bd = random.uniform(0.02, 0.04), random.uniform(0.18, 0.3), random.uniform(0.16, 0.22)
                if x + bw > x1 - 0.03:
                    break
                box(f"{prefix}_Book{k}", (x + bw / 2, yc, z + bh / 2), (bw, bd, bh), random.choice(M['books']))
                x += bw + 0.002
                k += 1
        elif r < 0.8:     # 花瓶
            vr, vh = random.uniform(0.04, 0.07), random.uniform(0.12, 0.26)
            if x + 2 * vr > x1 - 0.03:
                break
            cyl(f"{prefix}_Vase{k}", (x + vr, yc, z + vh / 2), vr, vh,
                random.choice((M['ceramic_white'], M['ceramic_blue'])), r2=vr * 0.6)
            x += 2 * vr + 0.04
            k += 1
        else:             # 斜靠的盤子
            pr = random.uniform(0.08, 0.13)
            if x + 2 * pr > x1 - 0.03:
                break
            cyl(f"{prefix}_Plate{k}", (x + pr, yc + 0.06, z + pr), pr, 0.015,
                random.choice((M['ceramic_white'], M['ceramic_blue'])), rot=(radians(75), 0, 0), seg=32)
            x += 2 * pr + 0.03
            k += 1
        x += 0.03


def build_bookshelf(M):
    use_coll("Bookshelf")
    tr = M['trim']
    x0, x1, d, top = 4.1, 6.3, 0.38, 2.45
    y0, cy, cx, t = D - d, D - d / 2, (x0 + x1) / 2, 0.035
    box("Shelf_Back", (cx, D - 0.01, H / 2), (x1 - x0, 0.02, H), M['wall'])
    for k, x in enumerate((x0 + t / 2, cx, x1 - t / 2)):
        box(f"Shelf_Side{k}", (x, cy, top / 2), (t, d, top), tr)
    box("Shelf_Cabinet", (cx, cy, 0.43), (x1 - x0, d, 0.86), tr)
    box("Shelf_Counter", (cx, cy - 0.02, 0.88), (x1 - x0 + 0.06, d + 0.04, 0.04), tr, bevel=0.008)
    nd = 4
    dw = (x1 - x0) / nd
    for i in range(nd):
        a = x0 + i * dw
        molding_rect(f"Shelf_Door{i}", None, a + 0.07, a + dw - 0.07, 0.14, 0.76, y0 - 0.008, tr, t=0.022, d=0.016)
        sphere(f"Shelf_Knob{i}", (a + dw - 0.06 if i % 2 == 0 else a + 0.06, y0 - 0.025, 0.55), 0.015, M['brass'])
    levels = [0.9, 1.29, 1.68, 2.07]
    for k, z in enumerate(levels[1:] + [top - 0.02]):
        box(f"Shelf_Board{k}", (cx, cy, z), (x1 - x0, d - 0.02, 0.03), tr)
    box("Shelf_Crown", (cx, cy - 0.03, top + 0.05), (x1 - x0 + 0.12, d + 0.06, 0.1), tr, bevel=0.01)
    box("Shelf_Header", (cx, cy, (top + H) / 2), (x1 - x0, d, H - top), tr)
    for k, z in enumerate(levels):
        for bay, (a, b) in enumerate(((x0 + t, cx - t / 2), (cx + t / 2, x1 - t))):
            fill_shelf(f"Shelf_{k}_{bay}", a, b, z + 0.015, cy + 0.02, M)


# ================================================================ 家具
def turned_leg(name, parent, x, y, h, mat_, s=1.0):
    """車削桌腳"""
    cyl(name + "_Foot", (x, y, 0.03 * h), 0.032 * s, 0.06 * h, mat_, parent)
    cyl(name + "_Shaft", (x, y, 0.28 * h), 0.022 * s, 0.44 * h, mat_, parent, r2=0.032 * s)
    sphere(name + "_Bulb", (x, y, 0.58 * h), 0.042 * s, mat_, parent, scale=(1, 1, 1.5))
    cyl(name + "_Ring", (x, y, 0.69 * h), 0.04 * s, 0.03 * h, mat_, parent)
    cyl(name + "_Neck", (x, y, 0.76 * h), 0.026 * s, 0.12 * h, mat_, parent)
    box(name + "_Block", (x, y, 0.91 * h), (0.07 * s, 0.07 * s, 0.18 * h), mat_, parent)


CORNERS = ((-1, -1), (1, -1), (-1, 1), (1, 1))


def build_sofa(name, loc, rot_z, M, L=2.2, Dp=0.95):
    e = empty(name, loc, rot_z)
    f, wd = M['sofa'], M['wood_dark']
    aw, lh, sz = 0.24, 0.12, 0.42
    iw = L - 2 * aw
    cw = iw / 3
    box(name + "_Base", (0, 0, (lh + sz) / 2), (L - 0.04, Dp - 0.04, sz - lh), f, e, bevel=0.03)
    box(name + "_Back", (0, Dp / 2 - 0.11, (lh + 0.88) / 2), (iw + 0.02, 0.22, 0.88 - lh), f, e, bevel=0.04)
    for i in range(3):
        x = -iw / 2 + cw * (i + 0.5)
        box(f"{name}_SeatCush{i}", (x, -0.11, sz + 0.08), (cw - 0.015, Dp - 0.24, 0.16), f, e, bevel=0.05)
        box(f"{name}_BackCush{i}", (x, Dp / 2 - 0.31, sz + 0.33), (cw - 0.02, 0.18, 0.36), f, e,
            rot=(-0.18, 0, 0), bevel=0.07)
    for s in (-1, 1):   # 捲手扶手
        x = s * (L / 2 - aw / 2)
        box(f"{name}_Arm{s}", (x, 0.0, (lh + 0.6) / 2), (aw, Dp, 0.6 - lh), f, e, bevel=0.03)
        cyl(f"{name}_ArmRoll{s}", (x + s * 0.01, 0.0, 0.6), aw / 2 + 0.025, Dp + 0.02, f, e,
            rot=(pi / 2, 0, 0), seg=32, bevel=0.02)
        sphere(f"{name}_Pillow{s}", (s * (iw / 2 - 0.25), Dp / 2 - 0.45, sz + 0.36), 0.22, M['pillow'], e,
               scale=(1, 0.35, 0.9), rot=(-0.25, 0, s * 0.25))
    for i, (sx, sy) in enumerate(CORNERS):
        cyl(f"{name}_Leg{i}", (sx * (L / 2 - 0.1), sy * (Dp / 2 - 0.1), lh / 2), 0.022, lh, wd, e, r2=0.035)


def build_wingchair(name, loc, rot_z, fab, M):
    e = empty(name, loc, rot_z)
    wd = M['wood_dark']
    Wd, Dp, lh, sz = 0.82, 0.82, 0.16, 0.42
    box(name + "_Base", (0, 0, (lh + sz) / 2), (Wd - 0.04, Dp - 0.04, sz - lh), fab, e, bevel=0.03)
    box(name + "_Cushion", (0, -0.06, sz + 0.06), (Wd - 0.3, Dp - 0.22, 0.12), fab, e, bevel=0.05)
    box(name + "_Back", (0, Dp / 2 - 0.09, (sz + 1.15) / 2), (Wd - 0.12, 0.16, 1.15 - sz), fab, e,
        rot=(-0.1, 0, 0), bevel=0.05)
    for s in (-1, 1):
        box(f"{name}_Wing{s}", (s * (Wd / 2 - 0.07), Dp / 2 - 0.24, 0.86), (0.1, 0.36, 0.52), fab, e,
            rot=(0, 0, s * 0.12), bevel=0.04)
        box(f"{name}_Arm{s}", (s * (Wd / 2 - 0.08), -0.04, (sz + 0.62) / 2), (0.14, Dp - 0.12, 0.62 - sz), fab, e,
            bevel=0.03)
        cyl(f"{name}_ArmRoll{s}", (s * (Wd / 2 - 0.07), -0.04, 0.62), 0.075, Dp - 0.1, fab, e,
            rot=(pi / 2, 0, 0), seg=24, bevel=0.015)
    for i, (sx, sy) in enumerate(CORNERS):
        rod(f"{name}_Leg{i}", (sx * (Wd / 2 - 0.08), sy * (Dp / 2 - 0.08), lh),
            (sx * (Wd / 2 - 0.05), sy * (Dp / 2 - 0.05), 0), 0.03, wd, e, r2=0.018)
    sphere(name + "_Pillow", (0, Dp / 2 - 0.25, sz + 0.32), 0.17, M['pillow'], e, scale=(1, 0.35, 0.85), rot=(-0.2, 0, 0))


def build_windsor(name, loc, rot_z, M):
    """溫莎椅"""
    e = empty(name, loc, rot_z)
    wd = M['wood_mid']
    sz = 0.45
    seat = cyl(name + "_Seat", (0, 0, sz), 0.22, 0.04, wd, e, seg=32, bevel=0.008)
    seat.scale = (1, 0.92, 1)
    cyl(name + "_Pad", (0, -0.01, sz + 0.03), 0.17, 0.02, M['cream'], e, seg=32, bevel=0.008)
    legs = []
    for i, (sx, sy) in enumerate(CORNERS):
        top, bot = Vector((sx * 0.13, sy * 0.12, sz - 0.02)), Vector((sx * 0.19, sy * 0.18, 0.0))
        rod(f"{name}_Leg{i}", top, bot, 0.022, wd, e, r2=0.016)
        legs.append((top, bot))

    def at(i, zz):
        t_, b_ = legs[i]
        return t_.lerp(b_, (t_.z - zz) / (t_.z - b_.z))

    zz = 0.16
    rod(name + "_StrL", at(0, zz), at(2, zz), 0.011, wd, e)
    rod(name + "_StrR", at(1, zz), at(3, zz), 0.011, wd, e)
    rod(name + "_StrM", (at(0, zz) + at(2, zz)) / 2, (at(1, zz) + at(3, zz)) / 2, 0.011, wd, e)
    by, R, zh = 0.16, 0.17, 0.8
    for s in (-1, 1):
        rod(f"{name}_Post{s}", (s * R, by, sz), (s * R, by + 0.03, zh), 0.014, wd, e)
    torus(name + "_Bow", (0, by + 0.03, zh), R, 0.016, wd, e, rot=(pi / 2, 0, 0), arc=pi, seg=24)
    for k, x in enumerate((-0.1, -0.05, 0.0, 0.05, 0.1)):
        rod(f"{name}_Spindle{k}", (x, by, sz), (x, by + 0.03, zh + sqrt(R * R - x * x) - 0.01), 0.007, wd, e, seg=6)


def build_coffee_table(name, loc, M):
    e = empty(name, loc)
    wd = M['wood_dark']
    L, Wd, Ht = 1.15, 0.62, 0.45
    box(name + "_Top", (0, 0, Ht - 0.02), (L, Wd, 0.04), wd, e, bevel=0.01)
    box(name + "_Apron", (0, 0, Ht - 0.08), (L - 0.1, Wd - 0.1, 0.08), wd, e)
    box(name + "_Shelf", (0, 0, 0.12), (L - 0.1, Wd - 0.1, 0.025), wd, e)
    for i, (sx, sy) in enumerate(CORNERS):
        turned_leg(f"{name}_Leg{i}", e, sx * (L / 2 - 0.07), sy * (Wd / 2 - 0.07), Ht - 0.04, wd, 0.9)
    box(name + "_Book1", (-0.3, 0, Ht + 0.02), (0.28, 0.2, 0.04), M['books'][1], e)
    box(name + "_Book2", (-0.3, 0.01, Ht + 0.06), (0.24, 0.17, 0.04), M['books'][3], e, rot=(0, 0, 0.2))
    sphere(name + "_Bowl", (0.25, 0, Ht + 0.04), 0.13, M['ceramic_blue'], e, scale=(1, 1, 0.35))
    for i in range(4):
        sphere(f"{name}_Fruit{i}", (0.25 + 0.05 * cos(i * 1.6), 0.05 * sin(i * 1.6), Ht + 0.08), 0.035, M['fruit'], e)


def build_round_table(name, loc, M):
    e = empty(name, loc)
    wl = M['wood_light']
    Ht = 0.58
    cyl(name + "_Top", (0, 0, Ht - 0.015), 0.32, 0.03, wl, e, seg=40, bevel=0.005)
    cyl(name + "_Shelf", (0, 0, 0.18), 0.26, 0.02, wl, e, seg=40)
    for i in range(3):
        a = 2 * pi * i / 3
        rod(f"{name}_Leg{i}", (0.2 * cos(a), 0.2 * sin(a), Ht - 0.03), (0.27 * cos(a), 0.27 * sin(a), 0), 0.018, wl, e)
    potted_plant(name + "_Fern", (0, 0, Ht), M, parent=e, pot_r=0.09, pot_h=0.14, n=14,
                 leaf_len=0.15, leaf_w=0.035, height=0.2, pot='ceramic_white', elev=(-15, 60))


def build_end_table_lamp(name, loc, M):
    e = empty(name, loc)
    wm = M['wood_mid']
    Ht, s = 0.6, 0.46
    box(name + "_Top", (0, 0, Ht - 0.02), (s, s, 0.04), wm, e, bevel=0.008)
    box(name + "_Drawer", (0, 0, Ht - 0.12), (s - 0.04, s - 0.04, 0.16), wm, e)
    sphere(name + "_Knob", (0, -s / 2 + 0.005, Ht - 0.12), 0.015, M['brass'], e)
    box(name + "_Shelf", (0, 0, 0.12), (s - 0.04, s - 0.04, 0.02), wm, e)
    for i, (sx, sy) in enumerate(CORNERS):
        cyl(f"{name}_Leg{i}", (sx * (s / 2 - 0.04), sy * (s / 2 - 0.04), (Ht - 0.2) / 2), 0.016, Ht - 0.2, wm, e, r2=0.022)
    sphere(name + "_LampBase", (0, 0, Ht + 0.13), 0.11, M['ceramic_blue'], e, scale=(1, 1, 1.15))
    cyl(name + "_LampNeck", (0, 0, Ht + 0.3), 0.012, 0.16, M['brass'], e)
    cyl(name + "_LampShade", (0, 0, Ht + 0.45), 0.19, 0.24, M['shade'], e, r2=0.12, caps=False, seg=40)
    light(name + "_Light", (0, 0, Ht + 0.42), 'POINT', 25, parent=e)


def build_living(M):
    use_coll("Living")
    box("Rug_Border", (3.0, 3.75, 0.004), (3.2, 2.6, 0.008), M['rug_border'])
    box("Rug", (3.0, 3.75, 0.009), (2.95, 2.35, 0.006), M['rug'])
    build_sofa("Sofa", (3.0, 2.4, 0), facing(0, 1), M)
    build_wingchair("Armchair_Cream", (1.25, 3.2, 0), facing(1.8, 0.4), M['cream'], M)
    build_wingchair("Armchair_Yellow", (5.15, 3.6, 0), facing(-1.6, 0.3), M['yellow'], M)
    build_coffee_table("CoffeeTable", (3.0, 3.8, 0), M)
    build_round_table("SideTable_Round", (1.35, 1.95, 0), M)
    build_end_table_lamp("EndTable", (4.45, 2.45, 0), M)


def build_dining(M):
    use_coll("Dining")
    tx, ty, L, Wd, Ht = 7.9, 4.3, 1.8, 0.95, 0.76
    e = empty("DiningTable", (tx, ty, 0))
    wm = M['wood_mid']
    box("DT_Top", (0, 0, Ht - 0.025), (L, Wd, 0.05), wm, e, bevel=0.012)
    box("DT_Apron", (0, 0, Ht - 0.11), (L - 0.16, Wd - 0.16, 0.12), wm, e)
    for i, (sx, sy) in enumerate(CORNERS):
        turned_leg(f"DT_Leg{i}", e, sx * (L / 2 - 0.12), sy * (Wd / 2 - 0.12), Ht - 0.05, wm, 1.1)

    cyl("DT_Vase", (0, 0, Ht + 0.1), 0.06, 0.2, M['ceramic_white'], e, r2=0.045)
    for i in range(10):
        a, rr = random.uniform(0, 2 * pi), random.uniform(0, 0.09)
        sphere(f"DT_Flower{i}", (rr * cos(a), rr * sin(a), Ht + 0.24 + random.uniform(0, 0.08)), 0.035, M['flower'], e)
    for s in (-1, 1):
        cyl(f"DT_CandleStick{s}", (s * 0.45, 0, Ht + 0.09), 0.02, 0.18, M['brass'], e, r2=0.012)
        cyl(f"DT_Candle{s}", (s * 0.45, 0, Ht + 0.24), 0.012, 0.12, M['candle'], e)

    for i, dx in enumerate((-0.55, 0.0, 0.55)):
        build_windsor(f"Chair_S{i}", (tx + dx, ty - 0.62, 0), facing(0, 1), M)
        build_windsor(f"Chair_N{i}", (tx + dx, ty + 0.62, 0), facing(0, -1), M)
        for sy in (-1, 1):
            cyl(f"DT_Plate{i}_{sy}", (dx, sy * 0.3, Ht + 0.006), 0.12, 0.012, M['ceramic_white'], e, seg=32)
    build_windsor("Chair_W", (tx - 1.15, ty, 0), facing(1, 0), M)
    build_windsor("Chair_E", (tx + 1.15, ty, 0), facing(-1, 0), M)


def build_kitchen(M, right):
    use_coll("Kitchen")
    tr = M['trim']
    x0, x1, y0, y1 = 7.3, 9.2, 2.0, 2.6
    cx, cy = (x0 + x1) / 2, (y0 + y1) / 2
    box("Counter_Body", (cx, cy, 0.46), (x1 - x0, y1 - y0, 0.84), tr)
    box("Counter_Toe", (cx, cy + 0.03, 0.04), (x1 - x0 - 0.06, y1 - y0 - 0.06, 0.08), M['wood_dark'])
    box("Counter_Top", (cx, cy, 0.9), (x1 - x0 + 0.06, y1 - y0 + 0.06, 0.04), M['wood_mid'], bevel=0.008)
    nd = 3
    dw = (x1 - x0) / nd
    for i in range(nd):
        a = x0 + i * dw
        molding_rect(f"Counter_Door{i}", None, a + 0.08, a + dw - 0.08, 0.15, 0.78, y0 - 0.008, tr)
        sphere(f"Counter_Knob{i}", (a + dw / 2, y0 - 0.025, 0.7), 0.015, M['brass'])
    cyl("Counter_Bowl", (8.7, 2.3, 0.96), 0.13, 0.08, M['ceramic_white'], r2=0.16, seg=32)
    for i in range(5):
        sphere(f"Counter_Fruit{i}", (8.7 + 0.06 * cos(i * 1.3), 2.3 + 0.06 * sin(i * 1.3), 1.02), 0.04, M['fruit'])
    cyl("Counter_Jar", (7.6, 2.4, 1.02), 0.06, 0.2, M['glass'])

    # 右牆吊櫃 (牆本地座標, 世界 y≈2.95)
    ux = D - 2.95
    box("UpperCab_Body", (ux, -0.17, 1.98), (0.75, 0.34, 0.76), tr, right)
    box("UpperCab_Crown", (ux, -0.19, 2.39), (0.85, 0.38, 0.06), tr, right, bevel=0.008)
    for k in (-1, 1):
        gx = ux + k * 0.18
        box(f"UpperCab_Glass{k}", (gx, -0.342, 2.0), (0.3, 0.004, 0.6), M['glass'], right)
        molding_rect(f"UpperCab_Door{k}", right, gx - 0.16, gx + 0.16, 1.68, 2.32, -0.348, tr, t=0.04, d=0.016)
    box("GlassRack", (ux, -0.17, 1.585), (0.7, 0.25, 0.02), M['wood_dark'], right)
    for i in range(5):
        gx = ux - 0.28 + i * 0.14
        cyl(f"WineGlass{i}_Foot", (gx, -0.17, 1.57), 0.035, 0.005, M['glass'], right)
        cyl(f"WineGlass{i}_Stem", (gx, -0.17, 1.53), 0.004, 0.08, M['glass'], right)
        cyl(f"WineGlass{i}_Bowl", (gx, -0.17, 1.44), 0.04, 0.1, M['glass'], right, r2=0.02, caps=False)

    # 吧台吊燈
    use_coll("Lighting_Fixtures")
    for k, x in enumerate((7.8, 8.7)):
        rod(f"Pendant{k}_Cord", (x, 2.3, H), (x, 2.3, 1.86), 0.004, M['iron'])
        cyl(f"Pendant{k}_Cap", (x, 2.3, 1.87), 0.04, 0.03, M['iron'])
        cyl(f"Pendant{k}_Shade", (x, 2.3, 1.78), 0.15, 0.16, M['shade'], r2=0.05, caps=False, seg=32)
        light(f"Pendant{k}_Light", (x, 2.3, 1.78), 'POINT', 25)


# ================================================================ 吊燈
def chandelier_shades(name, x, y, zc, n, R, M):
    """客廳: 鐵件 + 燈罩式吊燈"""
    use_coll("Lighting_Fixtures")
    e = empty(name, (x, y, zc))
    iron, brass = M['iron'], M['brass']
    drop = H - zc
    cyl(name + "_Rod", (0, 0, (drop + 0.15) / 2), 0.007, drop - 0.15, iron, e)
    cyl(name + "_Canopy", (0, 0, drop - 0.02), 0.07, 0.04, iron, e)
    cyl(name + "_Column", (0, 0, 0.05), 0.02, 0.25, iron, e)
    sphere(name + "_Ball", (0, 0, -0.12), 0.045, brass, e)
    cyl(name + "_Finial", (0, 0, -0.19), 0.002, 0.08, iron, e, r2=0.02)
    for i in range(n):
        a = 2 * pi * i / n
        ca, sa = cos(a), sin(a)
        p0 = (0.03 * ca, 0.03 * sa, -0.08)
        p1 = (0.5 * R * ca, 0.5 * R * sa, -0.16)
        p2 = (R * ca, R * sa, -0.08)
        p3 = (R * ca, R * sa, 0.0)
        rod(f"{name}_ArmA{i}", p0, p1, 0.007, iron, e, seg=8)
        rod(f"{name}_ArmB{i}", p1, p2, 0.007, iron, e, seg=8)
        rod(f"{name}_ArmC{i}", p2, p3, 0.007, iron, e, seg=8)
        cyl(f"{name}_Cup{i}", (R * ca, R * sa, 0.01), 0.03, 0.02, brass, e)
        cyl(f"{name}_Candle{i}", (R * ca, R * sa, 0.06), 0.011, 0.09, M['candle'], e)
        cyl(f"{name}_Shade{i}", (R * ca, R * sa, 0.12), 0.085, 0.13, M['shade'], e, r2=0.05, caps=False, seg=28)
        light(f"{name}_Light{i}", (R * ca, R * sa, 0.1), 'POINT', 12, parent=e)


def chandelier_ring(name, x, y, zc, n, R, M):
    """餐廳: 鐵件圓環燭台吊燈"""
    use_coll("Lighting_Fixtures")
    e = empty(name, (x, y, zc))
    iron = M['iron']
    drop = H - zc
    cyl(name + "_Rod", (0, 0, (drop + 0.35) / 2), 0.008, drop - 0.35, iron, e)
    cyl(name + "_Canopy", (0, 0, drop - 0.02), 0.07, 0.04, iron, e)
    sphere(name + "_Hub", (0, 0, 0.35), 0.03, iron, e)
    for k in range(4):
        a = 2 * pi * k / 4 + pi / 4
        rod(f"{name}_Chain{k}", (0, 0, 0.35), (R * cos(a), R * sin(a), 0), 0.005, iron, e, seg=6)
    torus(name + "_Ring", (0, 0, 0), R, 0.014, iron, e, seg=48)
    torus(name + "_RingInner", (0, 0, -0.06), R * 0.55, 0.008, iron, e, seg=32)
    for i in range(n):
        a = 2 * pi * i / n
        p = (R * cos(a), R * sin(a))
        cyl(f"{name}_Cup{i}", (*p, 0.02), 0.028, 0.02, iron, e)
        cyl(f"{name}_Candle{i}", (*p, 0.08), 0.012, 0.1, M['candle'], e)
        sphere(f"{name}_Flame{i}", (*p, 0.145), 0.012, M['flame'], e, scale=(1, 1, 1.7), seg=10)
        light(f"{name}_Light{i}", (*p, 0.16), 'POINT', 8, (1.0, 0.75, 0.45), parent=e)
        if i % 2 == 0:
            rod(f"{name}_Spoke{i}", (R * 0.55 * cos(a), R * 0.55 * sin(a), -0.06), (*p, 0), 0.005, iron, e, seg=6)


# ================================================================ 植物
def potted_plant(name, loc, M, parent=None, pot_r=0.2, pot_h=0.35, n=12, leaf_len=0.4, leaf_w=0.09,
                 height=0.5, pot='terracotta', elev=(20, 70)):
    e = empty(name, loc, random.uniform(0, 2 * pi), parent)
    cyl(name + "_Pot", (0, 0, pot_h / 2), pot_r * 0.72, pot_h, M[pot], e, r2=pot_r, seg=28)
    torus(name + "_Rim", (0, 0, pot_h), pot_r, 0.1 * pot_r, M[pot], e, seg=28, rseg=6)
    cyl(name + "_Soil", (0, 0, pot_h - 0.02), pot_r * 0.93, 0.01, M['soil'], e, seg=28)
    for i in range(n):
        yaw = random.uniform(0, 2 * pi)
        el = radians(random.uniform(*elev))
        sh = height * random.uniform(0.35, 1.0)
        base = Vector((random.uniform(-0.3, 0.3) * pot_r, random.uniform(-0.3, 0.3) * pot_r, pot_h - 0.02))
        tip = base + Vector((-sin(yaw) * 0.15 * sh, cos(yaw) * 0.15 * sh, sh))
        rod(f"{name}_Stem{i}", base, tip, 0.006, M['stem'], e, seg=6)
        L = leaf_len * random.uniform(0.7, 1.15)
        dv = Vector((-sin(yaw) * cos(el), cos(yaw) * cos(el), sin(el)))
        sphere(f"{name}_Leaf{i}", tip + dv * L, 1.0, random.choice((M['leaf'], M['leaf2'])), e,
               scale=(leaf_w * random.uniform(0.8, 1.2), L, 0.006), rot=(el, 0, yaw), seg=12)


def build_plants(M):
    use_coll("Plants")
    potted_plant("Plant_FrontLeft", (0.45, 0.85, 0), M, pot_r=0.25, pot_h=0.5, n=14, leaf_len=0.45, leaf_w=0.1, height=0.8)
    potted_plant("Plant_BackLeft", (0.55, 5.45, 0), M, pot_r=0.2, pot_h=0.4, n=12, leaf_len=0.4, leaf_w=0.09, height=0.7)
    potted_plant("Plant_DiningCorner", (9.6, 5.6, 0), M, pot_r=0.22, pot_h=0.45, n=12, leaf_len=0.42, height=0.8)
    potted_plant("Plant_FrenchDoor", (9.55, 0.6, 0), M, pot_r=0.2, pot_h=0.38, n=10, leaf_len=0.38, height=0.6,
                 pot='ceramic_blue')


# ================================================================ 場景 / 攝影機 / 渲染
def add_camera(name, loc, target, lens=16):
    cd = bpy.data.cameras.new(name)
    cd.lens = lens
    cd.clip_start = 0.05
    o = bpy.data.objects.new(name, cd)
    o.location = loc
    o.rotation_euler = (Vector(target) - Vector(loc)).to_track_quat('-Z', 'Y').to_euler()
    return _link(o)


def setup_scene():
    sc = bpy.context.scene
    sc.unit_settings.system = 'METRIC'

    use_coll("Cameras")
    cam = add_camera("Camera_Living", (5.3, 0.25, 1.5), (3.6, 6.0, 1.15), 16)     # 對應下方手繪圖
    add_camera("Camera_Dining", (8.9, 0.3, 1.6), (3.0, 4.8, 1.0), 18)           # 從廚房看向客廳
    sc.camera = cam

    world = sc.world or bpy.data.worlds.new("World")
    sc.world = world
    try:
        world.use_nodes = True
    except Exception:
        pass
    nt = world.node_tree
    bg = next((n for n in nt.nodes if n.type == 'BACKGROUND'), None)
    if bg is None:
        bg = nt.nodes.new("ShaderNodeBackground")
        out = nt.nodes.new("ShaderNodeOutputWorld")
        nt.links.new(bg.outputs[0], out.inputs[0])
    bg.inputs[0].default_value = (0.7, 0.82, 1.0, 1)
    bg.inputs[1].default_value = 1.2

    use_coll("Lights")
    light("Sun", (W / 2, D + 5, 8), 'SUN', 3.5, (1.0, 0.95, 0.85), rot=(radians(-50), 0, radians(25)))

    sc.render.engine = 'CYCLES'
    try:
        sc.cycles.samples = SAMPLES
        sc.cycles.use_denoising = True
    except Exception:
        pass
    sc.render.resolution_x, sc.render.resolution_y = 1920, 1080
    for vt in ('AgX', 'Filmic'):
        try:
            sc.view_settings.view_transform = vt
            break
        except TypeError:
            continue
    sc.view_settings.exposure = 0.5

    if SKETCH_LINES:
        try:
            sc.render.use_freestyle = True
            sc.render.line_thickness_mode = 'ABSOLUTE'
            sc.render.line_thickness = 1.0
            vl = bpy.context.view_layer
            vl.use_freestyle = True
            fs = vl.freestyle_settings
            fs.crease_angle = radians(135)
            ls = fs.linesets[0] if len(fs.linesets) else fs.linesets.new("Sketch")
            ls.select_by_visibility = True
            ls.select_silhouette = ls.select_border = ls.select_crease = True
            ls.linestyle.color = (0.18, 0.14, 0.1)
            ls.linestyle.thickness = 1.3
        except Exception as ex:
            print("Freestyle 無法啟用:", ex)


def clear_scene():
    obj = getattr(bpy.context, "object", None)
    if obj is not None and obj.mode != 'OBJECT':
        bpy.ops.object.mode_set(mode='OBJECT')
    for o in list(bpy.data.objects):
        bpy.data.objects.remove(o, do_unlink=True)
    for c in list(bpy.data.collections):
        bpy.data.collections.remove(c)
    for blocks in (bpy.data.meshes, bpy.data.materials, bpy.data.lights, bpy.data.cameras):
        for b in list(blocks):
            if b.users == 0:
                blocks.remove(b)


def main():
    t0 = time.time()
    clear_scene()
    M = build_materials()
    left, right = build_room(M)
    build_ceiling(M)
    build_tv_wall(M)
    build_bookshelf(M)
    build_living(M)
    build_dining(M)
    build_kitchen(M, right)
    chandelier_shades("Chandelier_Living", 3.0, 3.6, 2.2, 6, 0.38, M)
    chandelier_ring("Chandelier_Dining", 7.9, 4.3, 2.15, 8, 0.4, M)
    build_plants(M)
    setup_scene()
    print(f"[American Country Room] 完成: {len(bpy.data.objects)} 個物件, {time.time() - t0:.1f}s")


if __name__ == "__main__":
    main()
